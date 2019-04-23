package main

import (
	"fmt"
	"os"
	"strings"

	hssm "github.com/codacy/helm-ssm/internal"
	"github.com/spf13/cobra"
)

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

	cmd.MarkFlagRequired("values")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	funcMap := hssm.GetFuncMap()
	for _, filePath := range valueFiles {
		if err := hssm.ExecuteTemplate(filePath, funcMap, verbose, dryRun); err != nil {
			return err
		}
	}
	return nil
}

var valueFiles valueFilesList
var verbose bool
var dryRun bool

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
