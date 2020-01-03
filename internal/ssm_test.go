package hssm

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var fakeValue string = "value"
var fakeOtherValue string = "other-value"
var fakeMissingValue string = "missing-value"

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

		value, err := getSSMParameter(mockSvc, k, v.defaultValue, false)

		if v.expectedValue != nil && value != nil && *value == *v.expectedValue {
			// Success when expectedValue and value are both defined
			// and their values are equal
		} else if v.expectedValue == nil && v.value == nil && err != nil {
			// Success when expectedValue and value are both nil
			// getSSMParameter should return an error
		} else if err != nil {
			t.Errorf("Expected %s , got %s", *v.expectedValue, err)
		} else if value != nil {
			t.Errorf("Expected %s , got %s", *v.expectedValue, *value)
		} else {
			t.Errorf("Expected %s , got nil", *v.expectedValue)
		}
	}
}

func TestGetSSMParameterInvalidChar(t *testing.T) {
	key := "&%&/root/parameter5!$%&$&"
	t.Logf("Key with invalid characters should be handled")
	// Setup Test
	mockSvc := &mockSSMClient{}
	_, err := getSSMParameter(mockSvc, key, nil, false)
	if err == nil {
		t.Error(err)
	}
}

// GetParameter is a mock for the SSM client
func (m *mockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	parameterArn := "arn:::::"
	parameterLastModifiedDate := time.Now()
	parameterType := "String"
	parameterValue := fakeSSMStore[*input.Name]
	var parameterVersion int64 = 1

	if parameterValue.value == nil {
		return nil, awserr.New("ParameterNotFound", "", nil)
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
