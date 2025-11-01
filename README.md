# go-sdk

A Go-based SDK for processing templates, customizable with a YAML configuration.

## Overview

This library provides a flexible pipeline architecture for processing templates
with data and collecting results. It's designed for use in scaffolding tools,
code generators, and other systems that need to generate files from templates
with dynamic data.

Key features:

- Template processing pipeline with configurable data preprocessing;
- Support for multiple template providers (filesystem, etc.);
- Multiple collector types for different output strategies;
- Configurable filtering system;
- Prefix support for metadata and values.

## Architecture

The library consists of the following main components:

- **[Pipeline](./pkg/pipeline)**: the core processing engine that orchestrates
  template processing
- **[Template Providers](./pkg/templateproviders)**: sources for templates
  (e.g., filesystem, embedded)
- **[Collectors](./pkg/collectors)**: destinations for processed templates
  (e.g., filesystem writer, filters)
- **[Filters](./pkg/filters)**: mechanisms to control which templates files are
  processed
- **[Templates](./pkg/templates)**: utilities for processing Go templates
- **[Values Loaders](./pkg/values)**: utilities for loading data (Manifest and
  configurations).

## Installation

To use this library in your Go project, add it to your dependencies:

```bash
go get github.com/go-scaffold/go-sdk/v2
```

To use a specific version of the library, you can specify it in your go.mod
file:

```bash
# To use the latest version
go get github.com/go-scaffold/go-sdk/v2@latest

# To use a specific version (e.g., v2.1.0)
go get github.com/go-scaffold/go-sdk/v2@v2.0.1
```

After running these commands, your `go.mod` file will be updated to use the
specified version.

## Usage

### Basic Pipeline Usage

```go
package main

import (
  "text/template"
  "strings"
  "time"

  "github.com/go-scaffold/go-sdk/v2/pkg/collectors"
  "github.com/go-scaffold/go-sdk/v2/pkg/filters"
  "github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
  "github.com/go-scaffold/go-sdk/v2/pkg/templateproviders"
)

func main() {
  // Create a template provider (reads templates from filesystem)
  templateProvider := templateproviders.NewFileSystemProvider("./templates",
    filters.NewPatternFilter([]string{"**/*.tmpl", "**/*.gotmpl"}))

  // Create a collector (writes processed templates to filesystem)
  collector := collectors.NewFileWriterCollector("./output", nil)

  // Define template functions
  funcs := template.FuncMap{
    // Add your custom template functions here
    "upper": strings.ToUpper,
  }

  // Build the pipeline
  pipe, err := pipeline.NewPipelineBuilder().
    WithTemplateProvider(templateProvider).
    WithCollector(collector).
    WithFunctions(funcs).
    Build()
  if err != nil {
    panic(err)
  }

  // Define your data
  data := map[string]interface{}{
    "projectName": "MyProject",
    "version": "1.0.0",
  }

  metadata := map[string]interface{}{
    "generatedAt": time.Now().Unix(),
  }

  // Combine metadata and data, or use one set of data
  processData := map[string]interface{}{
    "metadata": metadata,
    "data":     data,
  }

  // Or merge the maps together if you prefer flat structure
  // processData := make(map[string]interface{})
  // for k, v := range metadata { processData[k] = v }
  // for k, v := range data { processData[k] = v }

  // Process all templates
  err = pipe.Process(processData)
  if err != nil {
    panic(err)
  }
}
```

### Using Collectors Chain

You can chain collectors to process templates in multiple ways:

```go
// Create a chain of collectors
filterCollector := collectors.NewFilterCollector(
  filters.NewPatternFilter([]string{"**/README.md", "**/main.go"}),
  collectors.NewFileWriterCollector("./filtered-output", nil),
)
fileWriter := collectors.NewFileWriterCollector("./all-output", filterCollector)

// Use in the pipeline
pipe, err := pipeline.NewPipelineBuilder().
  WithTemplateProvider(templateProvider).
  WithCollector(fileWriter).  // Processes all, passes matches to filterCollector
  WithFunctions(funcs).
  Build()
```

### Custom Data Preprocessing

You can preprocess your data before it's used in templates:

```go
preprocessor := func(data map[string]interface{}) (map[string]interface{}, error) {
  // Transform your data as needed
  data["computedValue"] = calculateSomething(data)
  return data, nil
}

pipe, err := pipeline.NewPipelineBuilder().
  WithTemplateProvider(templateProvider).
  WithCollector(collector).
  WithFunctions(funcs).
  WithDataPreprocessor(preprocessor).
  Build()
```

### Using Prefixes in Data

The pipeline doesn't have built-in prefix support, but you can achieve the same
result by wrapping your data manually before passing it to the pipeline:

```go
// Define your data and wrap it with prefixes to access in templates
data := map[string]interface{}{
  "projectName": "MyProject",
  "version": "1.0.0",
}

metadata := map[string]interface{}{
  "generatedAt": time.Now().Unix(),
}

// Wrap data with prefixes manually for use in templates
prefixedData := map[string]interface{}{
  "Values":   data,     // Access via .Values.key in templates
  "Manifest": metadata, // Access via .Manifest.key in templates
}

// Process all templates with prefixed data
err = pipe.Process(prefixedData)
if err != nil {
  panic(err)
}
```

In your templates, you can then use:

```go
{{.Values.projectName}}
{{.Manifest.generatedAt}}
```

### Using the Configurable Loader

The library provides a Loader type that allows you to customize the prefixes and
basenames used for manifest and values files:

```go
package main

import (
  "fmt"

  "github.com/go-scaffold/go-sdk/v2/pkg/values"
)

func main() {
  // Create a loader with default settings
  // Default: manifest prefix="Manifest", basename="Manifest", values prefix="Values", basename="values"
  defaultLoader := values.NewLoader()

  // Load YAML files using the loader
  // This will look for Manifest.yaml/Manifest.yml and values.yaml/values.yml in the specified directory
  data, err := defaultLoader.LoadYAMLs("./my-template", []string{})
  if err != nil {
    panic(err)
  }
  fmt.Printf("Loaded data: %+v\n", data)

  // Create a loader with custom settings
  customLoader := values.NewLoaderWithValues("MyManifest", "config", "MyValues", "data")

  // Load YAML files with custom prefixes and basenames
  // This will look for MyManifest.yaml/MyManifest.yml and data.yaml/data.yml in the specified directory
  customData, err := customLoader.LoadYAMLs("./my-custom-template", []string{})
  if err != nil {
    panic(err)
  }
  fmt.Printf("Custom loaded data: %+v\n", customData)
}
```

## Development

### Prerequisites

- Go 1.23 or higher
- GNU Make (for convenience targets)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/go-scaffold/go-sdk
cd go-sdk
```

2. Install dependencies:
```bash
go mod download
```

### Build Commands

To build the project:
```bash
make build
```

To clean build artifacts:
```bash
make clean
```

### Testing

To run tests:
```bash
make go-test
```

To run comprehensive checks (tests, linting, coverage):
```bash
make go-check
```

To generate a code coverage report:
```bash
make go-coverage
```

### Linting

To lint the Go files:
```bash
make go-lint
```

### Available Make Targets

You can see all available targets by running:
```bash
make help
```

Common targets include:
- `help` - Show help information
- `build` - Build the project
- `clean` - Remove build artifacts
- `go-check` - Run comprehensive checks
- `go-test` - Run unit tests
- `go-lint` - Lint Go files
- `go-coverage` - Generate coverage reports
- `go-dep-clean` - Clean unused dependencies
- `go-dep-download` - Download dependencies
- `go-dep-upgrade` - Upgrade dependencies

## Releases

This library follows semantic versioning. You can find the available releases on
the [GitHub releases page](https://github.com/go-scaffold/go-sdk/releases). Each
release is tagged with a version number in the format `vX.Y.Z`.

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE)
file.
