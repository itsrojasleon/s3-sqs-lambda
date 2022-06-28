package main

import (
	"log"
	"os"
	"strings"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambdaEvents "github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambdanodejs"
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

	lambdaArn := []string{"arn:aws:lambda:", cdkRegion, ":", cdkAccount, ":", "function:process"}

	// role := iam.NewRole(stack, jsii.String("CoolRoleForLambda"), &iam.RoleProps{
	// 	AssumedBy: iam.NewCompositePrincipal(
	// 		iam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &iam.ServicePrincipalOpts{}),
	// 		iam.NewServicePrincipal(jsii.String("sqs.amazonaws.com"), &iam.ServicePrincipalOpts{}),
	// 	),
	// 	ManagedPolicies: &[]iam.IManagedPolicy{
	// 		iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")),
	// 		iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaSQSQueueExecutionRole")),
	// 	},
	// 	InlinePolicies: &map[string]iam.PolicyDocument{
	// 		"hahaha": iam.NewPolicyDocument(&iam.PolicyDocumentProps{
	// 			Statements: &[]iam.PolicyStatement{
	// 				iam.NewPolicyStatement(&iam.PolicyStatementProps{
	// 					Effect: iam.Effect_ALLOW,
	// 					Actions: &[]*string{
	// 						jsii.String("sqs:ReceiveMessage"),
	// 					},
	// 					Resources: &[]*string{
	// 						jsii.String("*"),
	// 					},
	// 				}),
	// 			},
	// 		}),
	// 	},
	// })

	// role.AddManagedPolicy(
	// 	iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")),
	// )
	// role.AddManagedPolicy(
	// 	iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaSQSQueueExecutionRole")),
	// )

	// fn := lambda.NodejsFunction_FromFunctionAttributes(stack, jsii.String("CoolFunction"), &awslambda.FunctionAttributes{
	// 	FunctionArn: jsii.String(strings.Join(lambdaArn, "")),
	// Role:        role,
	// })

	fn := lambda.NodejsFunction_FromFunctionArn(
		stack,
		jsii.String("CoolLambda"),
		jsii.String(strings.Join(lambdaArn, "")),
	)

	invokeEventSource := lambdaEvents.NewSqsEventSource(queue, &lambdaEvents.SqsEventSourceProps{})

	fn.AddEventSource(invokeEventSource)

	cdk.NewCfnOutput(stack, jsii.String("lambda-name"), &cdk.CfnOutputProps{
		Value: jsii.String(*fn.FunctionName()),
	})

	return stack
}

func main() {
	app := cdk.NewApp(nil)
	if cdkAccount == "" || cdkRegion == "" {
		log.Fatal("missing env variables")
	}

	NewInfraStack(app, "InfraStack", &InfraStackProps{
		cdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *cdk.Environment {
	return &cdk.Environment{
		Account: jsii.String(cdkAccount),
		Region:  jsii.String(cdkRegion),
	}
}
