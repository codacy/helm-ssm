package hssm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"k8s.io/helm/pkg/engine"
)

// ExecuteTemplate loads a template file, executes is against a given function map and writes the output
func ExecuteTemplate(sourceFilePath string, funcMap template.FuncMap, verbose bool, dryRun bool) error {
	fileContent, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}
	t := template.New("ssmtpl").Funcs(funcMap)
	if _, err := t.Parse(string(fileContent)); err != nil {
		return err
	}
	var buf bytes.Buffer
	vals := map[string]interface{}{}
	if err := t.Execute(&buf, vals); err != nil {
		return err
	}
	if verbose {
		fmt.Println(string(buf.Bytes()))
	}
	if !dryRun {
		ioutil.WriteFile(sourceFilePath, buf.Bytes(), 0777)
	}

	return nil
}

// GetFuncMap builds the relevant function map to helm_ssm
func GetFuncMap() template.FuncMap {
	e := engine.New()

	// Clone the func map because we are adding context-specific functions.
	var funcMap template.FuncMap = map[string]interface{}{}
	for k, v := range e.FuncMap {
		funcMap[k] = v
	}
	funcMap["ssm"] = func(ssmPath string, options ...string) (string, error) {
		var awsSession = NewAWSSession()
		return resolveSSMParameter(awsSession, ssmPath, options)
	}
	return funcMap
}

func resolveSSMParameter(awsSession *session.Session, ssmPath string, options []string) (string, error) {
	opts, err := handleOptions(options)
	if err != nil {
		return "", err
	}
	required, _ := strconv.ParseBool(opts["required"])
	var svc *ssm.SSM
	if region, exists := opts["region"]; exists {
		svc = ssm.New(awsSession, aws.NewConfig().WithRegion(region))
	} else {
		svc = ssm.New(awsSession)
	}

	res, err := GetSSMParameter(svc, opts["prefix"]+ssmPath, required, false)
	if err != nil {
		return "", err
	}
	return res, nil
}

func handleOptions(options []string) (map[string]string, error) {
	opts := map[string]string{}
	for _, o := range options {
		split := strings.Split(o, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("Invalid option: %s", o)
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
