package hssm

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// GetSSMParameter gets a parameter from the AWS Simple Systems Manager service.
func GetSSMParameter(svc ssmiface.SSMAPI, name string, defaultValue *string, decrypt bool) (*string, error) {
	regex := "([a-zA-Z0-9\\.\\-_/]*)"
	r, _ := regexp.Compile(regex)
	match := r.FindString(name)
	if match == "" {
		return nil, fmt.Errorf("there is an invalid character in the name of the parameter: %s. It should match %s", name, regex)
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
		if defaultValue != nil {
			return defaultValue, nil
		}
		return nil, err
	}
	if aerr != nil {
		return nil, err
	}
	return param.Parameter.Value, nil
}
