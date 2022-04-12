package main

import (
	"context"
	"log"

	"github.com/kristofferlind/terraform-provider-mssql/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// format terraform examples
//go:generate terraform fmt -recursive ./examples/

// generate documentation
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	version string = "development"
)

func main() {
	options := tfsdk.ServeOpts{
		Name: "registry.terraform.io/kristofferlind/mssql",
	}

	err := tfsdk.Serve(context.Background(), provider.New(version), options)

	if err != nil {
		log.Fatal(err.Error())
	}
}
