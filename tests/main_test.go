package testing

import (
	"testing"

	testaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const awsRegion = "ap-southeast-1"

type payload struct {
	Version    string   `json:"version"`
	ID         string   `json:"id"`
	DetailType string   `json:"detail-type"`
	Source     string   `json:"source"`
	Account    string   `json:"account"`
	Time       string   `json:"time"`
	Region     string   `json:"region"`
	Resources  []string `json:"resources"`
	Detail     struct {
		InstanceID     string `json:"instance-id"`
		InstanceAction string `json:"instance-action"`
	} `json:"detail"`
}

type detail struct {
	InstanceID     string `json:"instance-id"`
	InstanceAction string `json:"instance-action"`
}

func TestDefault(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples",
		NoColor:      true,
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}
	// Make sure testing infrastructure removed at last
	defer terraform.Destroy(t, terraformOptions)
	// Do terraform init and terraform apply --auto-approve
	terraform.InitAndApply(t, terraformOptions)
	// Get the function name
	functionName := terraform.Output(t, terraformOptions, "function_name")
	// Initialize the payload
	payload := &payload{
		Version:    "0",
		ID:         "12345678-1234-1234-1234-743977200366",
		DetailType: "EC2 Spot Instance Interruption Warning",
		Source:     "aws.ec2",
		Account:    "12345678",
		Time:       "2020-02-01T20:00:23Z",
		Region:     "ap-southeast-1",
		Resources:  []string{"arn:aws:ec2:ap-southeast-1:1234567:instance/i-123456789"},
		Detail: detail{
			InstanceID:     "i-123456789",
			InstanceAction: "terminate",
		},
	}
	// Invoke lambda function with sample payload
	testaws.InvokeFunction(t, awsRegion, functionName, *payload)
}
