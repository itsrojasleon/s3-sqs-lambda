package main

import (
	"fmt"
	"log"
	"os"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	lambdaEvents "github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	lambdaNode "github.com/aws/aws-cdk-go/awscdk/v2/awslambdanodejs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	s3n "github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	sqs "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	cdk.StackProps
}

var cdkAccount = os.Getenv("CDK_DEFAULT_ACCOUNT")
var cdkRegion = os.Getenv("CDK_DEFAULT_REGION")

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) cdk.Stack {
	var sprops cdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := cdk.NewStack(scope, &id, &sprops)

	bucket := s3.NewBucket(stack, jsii.String("CoolBucket"), &s3.BucketProps{
		Versioned: jsii.Bool(false),
	})

	queue := sqs.NewQueue(stack, jsii.String("CoolQueue"), &sqs.QueueProps{})

	bucket.AddEventNotification(s3.EventType_OBJECT_CREATED_PUT, s3n.NewSqsDestination(queue))

	lambdaArn := fmt.Sprintf(
		"%s%s%s%s%s%s",
		"arn:aws:lambda:",
		cdkRegion,
		":",
		cdkAccount,
		":",
		"function:process",
	)

	roleArn := fmt.Sprintf(
		"%s%s%s",
		"arn:aws:iam::",
		cdkAccount,
		":role/s3-sqs-lambda-dev-us-east-1-lambdaRole",
	)

	role := iam.Role_FromRoleArn(
		stack,
		jsii.String("LambdaRole"),
		jsii.String(roleArn),
		&iam.FromRoleArnOptions{
			Mutable: jsii.Bool(false),
		},
	)

	fn := lambdaNode.NodejsFunction_FromFunctionAttributes(
		stack,
		jsii.String("CoolLambda"),
		&lambda.FunctionAttributes{
			Role:        role,
			FunctionArn: jsii.String(lambdaArn),
		},
	)

	sqsEventSource := lambdaEvents.NewSqsEventSource(queue, &lambdaEvents.SqsEventSourceProps{
		// This will do the magic.
		MaxBatchingWindow: cdk.Duration_Seconds(jsii.Number(5)),
	})

	fn.AddEventSource(sqsEventSource)

	return stack
}

func main() {
	app := cdk.NewApp(nil)

	if cdkAccount == "" || cdkRegion == "" {
		log.Fatal("Missing environment variables")
	}

	NewInfraStack(app, "InfraStack", &InfraStackProps{
		cdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *cdk.Environment {
	return &cdk.Environment{
		Account: jsii.String(cdkAccount),
		Region:  jsii.String(cdkRegion),
	}
}
