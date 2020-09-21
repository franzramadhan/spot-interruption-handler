package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type eventDetail struct {
	InstanceID     string `json:"instance-id"`
	InstanceAction string `json:"instance-action"`
}

func handleRequest(ctx context.Context, event events.CloudWatchEvent) {
	// Create new session
	sess := session.Must(session.NewSession())

	// Declare variables
	var detail eventDetail

	// Validate event detail
	if strings.Contains(event.DetailType, "EC2 Spot Instance Interruption Warning") && event.Source == "aws.ec2" {
		// Parse event detail
		err := json.Unmarshal(event.Detail, &detail)
		if err != nil {
			log.Printf("ERROR - Unable to parse event detail - %v", err)
		}

		// Initialize autoscaling client
		asg := autoscaling.New(sess)
		// Get the autoscaling group metadata
		autoScalingInstanceOutput, err := asg.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{
			InstanceIds: []*string{
				aws.String(detail.InstanceID),
			},
		})
		if err != nil {
			log.Printf("ERROR - Unable to get autoscaling group metadata of of %s - %v", detail.InstanceID, err)
		}

		for _, asgInstance := range autoScalingInstanceOutput.AutoScalingInstances {
			log.Printf("INFO - Handling %s of %s", event.DetailType, aws.StringValue(asgInstance.AutoScalingGroupName))
			// If instance < 2, scale out to 2 to avoid unexpected downtime during detach operation
			if len(autoScalingInstanceOutput.AutoScalingInstances) < 2 {
				_, err = asg.SetDesiredCapacity(&autoscaling.SetDesiredCapacityInput{
					AutoScalingGroupName: asgInstance.AutoScalingGroupName,
					DesiredCapacity:      aws.Int64(2),
					HonorCooldown:        aws.Bool(false),
				})
				if err != nil {
					log.Printf("ERROR - Unable to increase capacity of autoscaling-group %s - %v", aws.StringValue(asgInstance.AutoScalingGroupName), err)
				}
				// Wait until scale-out operation is completed and desired capacity with InService status > 1
				_ = func(ctx aws.Context, input *autoscaling.DescribeAutoScalingGroupsInput) error {
					w := request.Waiter{
						Name:        "WaitUntilNoSingleInstance",
						MaxAttempts: 15,
						Delay:       request.ConstantWaiterDelay(4 * time.Second),
						Acceptors: []request.WaiterAcceptor{
							{
								State:   request.SuccessWaiterState,
								Matcher: request.PathWaiterMatch, Argument: "contains(AutoScalingGroups[].[length(Instances[?LifecycleState=='InService']) == DesiredCapacity ][], `false`)",
								Expected: true,
							},
							{
								State:   request.RetryWaiterState,
								Matcher: request.PathWaiterMatch, Argument: "contains(AutoScalingGroups[].[length(Instances[?LifecycleState=='InService']) < DesiredCapacity ][], `false`)",
								Expected: false,
							},
						},
						Logger: asg.Config.Logger,
						NewRequest: func(opts []request.Option) (*request.Request, error) {
							req, _ := asg.DescribeAutoScalingGroupsRequest(input)
							req.SetContext(ctx)
							return req, nil
						},
					}
					return w.WaitWithContext(ctx)
				}(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						asgInstance.AutoScalingGroupName,
					},
				})
			}

			// Detach the instance from autoscaling group to also drain the connection from Load Balancer
			_, err = asg.DetachInstances(&autoscaling.DetachInstancesInput{
				AutoScalingGroupName: asgInstance.AutoScalingGroupName,
				InstanceIds: []*string{
					aws.String(detail.InstanceID),
				},
				ShouldDecrementDesiredCapacity: aws.Bool(false),
			})
			if err != nil {
				log.Printf("ERROR - Unable to detach %s from autoscaling-group %s - %v", detail.InstanceID, aws.StringValue(asgInstance.AutoScalingGroupName), err)
			}
		}
	}
}

func main() {
	_, isLambda := os.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	if isLambda {
		lambda.Start(handleRequest)
	} else {
		log.Fatal("This is only intended to be run inside AWS Lambda")
	}
}
