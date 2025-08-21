package hssm

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var fakeValue = "value"
var fakeOtherValue = "other-value"
var fakeMissingValue = "missing-value"

type SSMParameter struct {
	value         *string
	defaultValue  *string
	expectedValue *string
}

var fakeSSMStore = map[string]SSMParameter{
	"/root/existing-parameter":                     SSMParameter{&fakeValue, nil, &fakeValue},
	"/root/existing-parameter-with-default":        SSMParameter{&fakeValue, &fakeOtherValue, &fakeValue},
	"/root/non-existing-parameter":                 SSMParameter{nil, &fakeMissingValue, &fakeMissingValue},
	"/root/non-existing-parameter-without-default": SSMParameter{nil, nil, nil},
}

type mockSSMClient struct {
	ssmiface.SSMAPI
}

func TestGetSSMParameter(t *testing.T) {
	// Setup Test
	mockSvc := &mockSSMClient{}

	for k, v := range fakeSSMStore {
		expectedValueStr := "nil"
		if v.expectedValue != nil {
			expectedValueStr = *v.expectedValue
		}
		t.Logf("Key: %s should have value: %s", k, expectedValueStr)

		value, err := GetSSMParameter(mockSvc, k, v.defaultValue, false)

		assert.Equal(t, v.expectedValue, value)
		if v.expectedValue == nil {
			assert.Error(t, err, "ParameterNotFound: Parameter does not exist in SSM")
		}
	}
}

func TestGetSSMParameterInvalidChar(t *testing.T) {
	key := "&%&/root/parameter5!$%&$&"
	t.Logf("Key with invalid characters should be handled")
	// Setup Test
	mockSvc := &mockSSMClient{}
	_, err := GetSSMParameter(mockSvc, key, nil, false)
	assert.Error(t, err, "there is an invalid character in the name of the parameter: &%&/root/parameter5!$%&$&. It should match ([a-zA-Z0-9\\.\\-_/]*)")
}

// GetParameter is a mock for the SSM client
func (m *mockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	parameterArn := "arn:::::"
	parameterLastModifiedDate := time.Now()
	parameterType := "String"
	parameterValue := fakeSSMStore[*input.Name]
	var parameterVersion int64 = 1

	if parameterValue.value == nil {
		return nil, awserr.New("ParameterNotFound", "Parameter does not exist in SSM", nil)
	}

	parameter := ssm.Parameter{
		ARN:              &parameterArn,
		LastModifiedDate: &parameterLastModifiedDate,
		Name:             input.Name,
		Type:             &parameterType,
		Value:            parameterValue.value,
		Version:          &parameterVersion,
	}
	getParameterOutput := &ssm.GetParameterOutput{
		Parameter: &parameter,
	}

	return getParameterOutput, nil
}
