package main

import (
	"flag"
	"fmt"
	"github.com/pCilip/GraphqlGoGenerate/internal/Generator"
	"github.com/pCilip/GraphqlGoGenerate/internal/SchemaProvider"
	"github.com/pCilip/GraphqlGoGenerate/internal/SchemaProvider/HttpProvider"
	"github.com/pCilip/GraphqlGoGenerate/internal/SchemaProvider/JsonProvider"
	"os"
)

func main() {
	packageName := flag.String("p", "", "Used to generate package paths for correct project.")
	filePath := flag.String("f", "", "File path used to load graphql schema.")
	httpEndpoint := flag.String("h", "", "Http endpoint where schema can be downloaded.")

	flag.Parse()

	if *packageName == "" {
		fmt.Println("p - not set")
		os.Exit(1)
	}

	if *filePath == "" && *httpEndpoint == "" {
		fmt.Println("f or h not set")
		os.Exit(1)
	}

	if *filePath != "" && *httpEndpoint != "" {
		fmt.Println("only one of f or h can be set")
		os.Exit(1)
	}

	var provider SchemaProvider.SchemaProvider

	if *filePath != "" {
		provider = &JsonProvider.Provider{FilePath: *filePath}
	}

	if *httpEndpoint != "" {
		provider = &HttpProvider.Provider{HttpEndpoint: *httpEndpoint}
	}

	if provider == nil {
		fmt.Println("Schema provider not set")
		os.Exit(1)
	}
	// load schema to memory using interface....
	schema := provider.MustProvide()
	gen := Generator.NewGenerator(schema.Data, *packageName)
	gen.Generate()
}
