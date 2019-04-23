package hssm

import (
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
func GetSSMParameter(svc ssmiface.SSMAPI, name string, required bool, decrypt bool) (string, error) {
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
			return "\"\"", nil
		}
		return "", err
	}
	if aerr != nil {
		return "", err
	}
	return *param.Parameter.Value, nil
}
