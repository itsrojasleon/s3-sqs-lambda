package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	bucket := awss3.NewBucket(stack, jsii.String("coolBucket"), &awss3.BucketProps{
		Versioned: jsii.Bool(false),
	})

	queue := awssqs.NewQueue(stack, jsii.String("coolQueue"), &awssqs.QueueProps{
		ReceiveMessageWaitTime: awscdk.Duration_Seconds(jsii.Number(5)),
	})

	role := awsiam.NewRole(stack, jsii.String("lambdaRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(
			jsii.String("lambda.amazonaws.com"),
			&awsiam.ServicePrincipalOpts{},
		),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaSQSQueueExecutionRole")),
		},
	})

	bucket.AddEventNotification(
		awss3.EventType_OBJECT_CREATED_PUT,
		awss3notifications.NewSqsDestination(queue),
	)

	awscdk.NewCfnOutput(stack, jsii.String("lambdaRoleArn"), &awscdk.CfnOutputProps{
		Value:      role.RoleArn(),
		ExportName: jsii.String("lambdaRoleArn"),
	})

	awscdk.NewCfnOutput(
		stack,
		jsii.String("coolQueueArn"),
		&awscdk.CfnOutputProps{
			Value:      queue.QueueArn(),
			ExportName: jsii.String("coolQueueArn"),
		})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewInfraStack(app, "InfraStack", &InfraStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
