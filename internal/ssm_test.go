package hssm

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type GetSSMParameterTestValue struct {
	path          string
	required      bool
	decrypt       bool
	expectedValue string
}

var ssmTestValues = map[string]GetSSMParameterTestValue{
	"/root/parameter1": GetSSMParameterTestValue{"/root/parameter1", false, false, "value1"},
	"/root/parameter2": GetSSMParameterTestValue{"/root/parameter2", false, false, "value2"},
	"/root/parameter3": GetSSMParameterTestValue{"/root/parameter3", false, false, "value3"},
	"/root/parameter4": GetSSMParameterTestValue{"/root/parameter4", false, false, "value4"},
	"/root/parameter5": GetSSMParameterTestValue{"/root/parameter5", false, false, "value5"},
}

type mockSSMClient struct {
	ssmiface.SSMAPI
}

func TestGetSSMParameter(t *testing.T) {
	// Setup Test
	mockSvc := &mockSSMClient{}

	for k, v := range ssmTestValues {
		t.Logf("Key: %s should have value: %s", k, v.expectedValue)
		value, err := getSSMParameter(mockSvc, v.path, v.required, v.decrypt)
		if err != nil {
			t.Errorf("Expected %s , got %s", v.expectedValue, err)
		} else if value != v.expectedValue {
			t.Errorf("Expected %s , got %s", v.expectedValue, value)
		}
	}
}

func TestGetSSMParameterInvalidChar(t *testing.T) {
	key := "&%&/root/parameter5!$%&$&"
	t.Logf("Key with invalid characters should be handled")
	// Setup Test
	mockSvc := &mockSSMClient{}
	_, err := getSSMParameter(mockSvc, key, false, false)
	if err == nil {
		t.Error(err)
	}
}

// GetParameter is a mock for the SSM client
func (m *mockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	parameterArn := "arn:::::"
	parameterLastModifiedDate := time.Now()
	parameterType := "String"
	parameterValue := ssmTestValues[*input.Name]
	var parameterVersion int64 = 1

	parameter := ssm.Parameter{
		ARN:              &parameterArn,
		LastModifiedDate: &parameterLastModifiedDate,
		Name:             input.Name,
		Type:             &parameterType,
		Value:            &parameterValue.expectedValue,
		Version:          &parameterVersion,
	}
	getParameterOutput := &ssm.GetParameterOutput{
		Parameter: &parameter,
	}
	return getParameterOutput, nil
}
