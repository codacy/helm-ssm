package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	hssm "github.com/tutti-ch/helm-ssm/internal"
)

var valueFiles valueFilesList
var targetDir string
var profile string
var verbose bool
var dryRun bool
var clean bool
var tagCleaned string
var prefix string

type valueFilesList []string

func (v *valueFilesList) String() string {
	return fmt.Sprint(*v)
}

func (v *valueFilesList) Type() string {
	return "valueFilesList"
}

func (v *valueFilesList) Set(value string) error {
	for _, filePath := range strings.Split(value, ",") {
		*v = append(*v, filePath)
	}
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:   "ssm [flags]",
		Short: "",
		RunE:  run,
	}

	f := cmd.Flags()
	f.VarP(&valueFiles, "values", "f", "specify values in a YAML file (can specify multiple)")
	f.BoolVarP(&verbose, "verbose", "v", false, "show the computed YAML values file/s")
	f.BoolVarP(&dryRun, "dry-run", "d", false, "doesn't replace the file content")
	f.StringVarP(&targetDir, "target-dir", "o", "", "dir to output content")
	f.StringVarP(&profile, "profile", "p", "", "aws profile to fetch the ssm parameters")
	f.BoolVarP(&clean, "clean", "c", false, "clean all template commands from file")
	f.StringVarP(&tagCleaned, "tag-cleaned", "t", "", "replace cleaned template commands with given string")
	f.StringVarP(&prefix, "prefix", "P", "", "prefix for all parameters without affecting the path. ignored if individual prefix is defined")

	err := cmd.MarkFlagRequired("values")
	if err != nil {
		fmt.Println("Error marking flag required:", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	funcMap := hssm.GetFuncMap(profile, prefix, clean, tagCleaned)
	for _, filePath := range valueFiles {
		content, err := hssm.ExecuteTemplate(filePath, funcMap, verbose)
		if err != nil {
			return err
		}
		if !dryRun {
			err := write(filePath, targetDir, content)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", filePath, err)
			}
		}
	}
	return nil
}

func write(filePath string, targetDir string, content string) error {
	if targetDir != "" {
		fileName := filepath.Base(filePath)
		return hssm.WriteFileD(fileName, targetDir, content)
	}
	return hssm.WriteFile(filePath, content)
}
