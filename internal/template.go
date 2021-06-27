package hssm

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/Masterminds/sprig"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"

)

var account_profiles = map[string]string{
	"staging":    "staging",
	"production": "default",
}

// WriteFileD dumps a given content on the file with path `targetDir/fileName`.
func WriteFileD(fileName string, targetDir string, content string) error {
	targetFilePath := targetDir + "/" + fileName
	_ = os.Mkdir(targetDir, os.ModePerm)
	return WriteFile(targetFilePath, content)
}

// WriteFile dumps a given content on the file with path `targetFilePath`.
func WriteFile(targetFilePath string, content string) error {
	return ioutil.WriteFile(targetFilePath, []byte(content), 0777)
}

// ExecuteTemplate loads a template file, executes is against a given function map and writes the output
func ExecuteTemplate(sourceFilePath string, funcMap template.FuncMap, verbose bool) (string, error) {
	fileContent, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return "", err
	}
	t := template.New("ssmtpl").Funcs(funcMap)
	if _, err := t.Parse(string(fileContent)); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	vals := map[string]interface{}{}
	if err := t.Execute(&buf, vals); err != nil {
		return "", err
	}
	if verbose {
		fmt.Println(string(buf.Bytes()))
	}
	return buf.String(), nil
}

// GetFuncMap builds the relevant function map to helm_ssm
func GetFuncMap(profile string) template.FuncMap {
	// Clone the func map because we are adding context-specific functions.
	var funcMap template.FuncMap = map[string]interface{}{}
	for k, v := range sprig.GenericFuncMap() {
		funcMap[k] = v
	}

	funcMap["ssm"] = func(ssmPath string, options ...string) (string, error) {
		optStr, err := resolveSSMParameter(ssmPath, options)
		str := ""
		if optStr != nil {
			str = *optStr
		}
		return str, err
	}
	return funcMap
}

func resolveSSMParameter(ssmPath string, options []string) (*string, error) {
	opts, err := handleOptions(options)
	if err != nil {
		return nil, err
	}

	var defaultValue *string
	if optDefaultValue, exists := opts["default"]; exists {
		defaultValue = &optDefaultValue
	}

	var session *session.Session

	_, lambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	if lambda {
		session = getLambdaSession(opts["account"])
	} else {
		session = getLocalSession(opts["account"])		
	}

	var svc ssmiface.SSMAPI
	if region, exists := opts["region"]; exists {
		svc = ssm.New(session, aws.NewConfig().WithRegion(region))
	} else {
		svc = ssm.New(session)
	}

	return GetSSMParameter(svc, opts["prefix"]+ssmPath, defaultValue, true)
}

func handleOptions(options []string) (map[string]string, error) {
	validOptions := []string{
		"required",
		"prefix",
		"region",
		"account",
	}
	opts := map[string]string{}
	for _, o := range options {
		split := strings.Split(o, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("Invalid option: %s. Valid options: %s", o, validOptions)
		}
		opts[split[0]] = split[1]
	}
	if _, exists := opts["required"]; !exists {
		opts["required"] = "true"
	}
	if _, exists := opts["prefix"]; !exists {
		opts["prefix"] = ""
	}
	return opts, nil
}

func getLocalSession(account string) *session.Session {
	s := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 account_profiles["account"],
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}))
	return s
}

func getLambdaSession(account string) *session.Session {
	var s *session.Session

	region := os.Getenv("AWS_REGION")
	s, _ = session.NewSession(&aws.Config{
		Region: &region,
	})
	if account != os.Getenv("PIPELINE_ENVIRONMENT") {
		cross_account_role_arn := os.Getenv("CROSS_ACCOUNT_ARN")
		s, _ = getAssumedSession(s, cross_account_role_arn, region)
	}
	return s
}

func getAssumedSession(baseSess *session.Session, roleArn, region string) (*session.Session, error) {
	stsSvc := sts.New(baseSess)
	assumedRole, _ := stsSvc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn: aws.String(roleArn),
	})

	return session.NewSession(&aws.Config{
		Credentials:   credentials.NewStaticCredentials(
			*assumedRole.Credentials.AccessKeyId,
			*assumedRole.Credentials.SecretAccessKey,
			*assumedRole.Credentials.SessionToken),
		Region:        aws.String(region),
	})
}
