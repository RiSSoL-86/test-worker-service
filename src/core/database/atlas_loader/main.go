package main

import (
	"app/src/core/models"
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(models.GetModels()...)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load Gorm schema: %v\n", err)
		os.Exit(1)
	}
	if _, err := io.WriteString(os.Stdout, stmts); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write Gorm schema: %v\n", err)
		os.Exit(1)
	}
}
