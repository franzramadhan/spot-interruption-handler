# Terraform AWS Lambda Spot Interruption Handler

This module will provision AWS Lambda function that subcribe to Cloudwatch Event to handle [EC2 Spot Interruption](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-interruptions.html).

This Lambda will detach EC2 spot instance from AutoScaling Group to drain the connection from ALB after Cloudwatch Event for EC2 Spot Interruption is triggered

## Table of Content

- [Terraform AWS Lambda Spot Interruption Handler](#terraform-aws-lambda-spot-interruption-handler)
  - [Table of Content](#table-of-content)
  - [Disclaimer](#disclaimer)
  - [Dependencies](#dependencies)
  - [Quick Start](#quick-start)
  - [Contributor](#contributor)
  - [License](#license)

## Disclaimer

- Spot instance interruptions notice happens 2 minutes prior the activity. So make sure your service does not contain long running process and can be exited safely during the lambda execution.
- For critical production workloads that cannot bear simultaneous spot interruptions, please use mixed instance distribution to ensure there is an on-demand instance for back-up `e.g: 50% spot and 50% on demand`.

## Dependencies

- [tfenv](https://github.com/tfutils/tfenv) for terraform binary management
- [go](https://golang.org/dl/) to run the test and build the lambda code
- [pre-commit-terraform](https://github.com/antonbabenko/pre-commit-terraform)
- [pre-commit-hooks](https://github.com/pre-commit/pre-commit-hooks)

## Quick Start

- Install [dependencies](#dependencies)
- Execute `pre-commit install`
- Go to `examples` and go to each scenario
- Follow instruction in `README.md`
- Go to `tests` to test using `terratest`.
  
```go
cd tests
go test -count 1 -v tests
```

## Contributor

For question, issue, and pull request you can contact these people:

- [Frans Caisar Ramadhan](https://github.com/franzramadhan) (**Author**)

## License

See the [LICENSE](https://github.com/franzramadhan/spot-interruption-handler/blob/master/LICENSE)
