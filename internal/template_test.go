package hssm

import (
	"io/ioutil"
	"syscall"
	"testing"
	"text/template"
)

func createTempFile() (string, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func TestExecuteTemplate(t *testing.T) {
	templateContent := "example: {{ and true false }}"
	expectedOutput := "example: false"
	t.Logf("Template with content: %s , should out put a file with content: %s", templateContent, expectedOutput)

	templateFilePath, err := createTempFile()
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(templateFilePath)
	ioutil.WriteFile(templateFilePath, []byte(templateContent), 0644)
	content, err := ExecuteTemplate(templateFilePath, template.FuncMap{}, false)
	if content != expectedOutput {
		t.Errorf("Expected content \"%s\". Got \"%s\"", expectedOutput, content)
	}
}

func TestWriteFile(t *testing.T) {
	templateContent := "write_file_example: true"
	expectedOutput := "write_file_example: true"
	t.Logf("Template with content: %s , should out put a file with content: %s", templateContent, expectedOutput)

	templateFilePath, err := createTempFile()
	if err != nil {
		panic(err)
	}
	WriteFile(templateFilePath, templateContent)
	fileContent, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		panic(err)
	}
	content := string(fileContent)
	if content != expectedOutput {
		t.Errorf("Expected file with content \"%s\". Got \"%s\"", expectedOutput, content)
	}
}
func TestFailExecuteTemplate(t *testing.T) {
	t.Logf("Non existing template should return \"no such file or directory\" error.")
	_, err := ExecuteTemplate("", template.FuncMap{}, false)
	if err == nil {
		t.Error("Should fail with \"no such file or directory\"")
	}
}

func TestGetFuncMap(t *testing.T) {
	t.Logf("\"ssm\" function should exist in function map.")
	funcMap := GetFuncMap()
	keys := make([]string, len(funcMap))
	for k := range funcMap {
		keys = append(keys, k)
	}
	if _, exists := funcMap["ssm"]; !exists {
		t.Errorf("Expected \"ssm\" function in function map. Got the following functions: %s", keys)
	}
}

func TestResolveSSMParameter(t *testing.T) {
	t.Logf("TODO")
}

func TestHandleOptions(t *testing.T) {
	t.Logf("TODO")
}
