package hssm

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// NewAWSSession loads a new session from shared config
func NewAWSSession() *session.Session {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return session
}

// GetSSMParameter gets a parameter from the AWS Simple Systems Manager service.
func GetSSMParameter(name string, required bool) (string, error) {
	return GetSSMParameterRD(name, required, false, "")
}

// GetSSMParameterR gets a parameter from the AWS Simple Systems Manager service, on a specific region.
func GetSSMParameterR(name string, required bool, region string) (string, error) {
	return GetSSMParameterRD(name, required, false, region)
}

// GetSSMParameterRD gets a parameter from the AWS Simple Systems Manager service.
func GetSSMParameterRD(name string, required bool, decrypt bool, region string) (string, error) {
	awsSession := NewAWSSession()
	if region != "" {
		svc := ssm.New(awsSession, aws.NewConfig().WithRegion(region))
		return getSSMParameter(svc, name, required, decrypt)
	}
	svc := ssm.New(awsSession)
	return getSSMParameter(svc, name, required, decrypt)
}

func getSSMParameter(svc ssmiface.SSMAPI, name string, required bool, decrypt bool) (string, error) {
	regex := "([a-zA-Z0-9\\.\\-_/]*)"
	r, _ := regexp.Compile(regex)
	match := r.FindString(name)
	if match == "" {
		return "", fmt.Errorf("There is an invalid character in the name of the parameter: %s. It should match %s", name, regex)
	}
	// Create the request to SSM
	getParameterInput := &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &decrypt,
	}

	// Get the parameter from SSM
	param, err := svc.GetParameter(getParameterInput)
	// Cast err to awserr.Error to handle specific error codes.
	aerr, ok := err.(awserr.Error)
	if ok && aerr.Code() == ssm.ErrCodeParameterNotFound {
		// Specific error code handling
		if !required {
			return "", nil
		}
		return "", err
	}
	if aerr != nil {
		return "", err
	}
	return *param.Parameter.Value, nil
}
