# AI Agent Context for go-sdk

## Project Overview
This is a Go-based template processing SDK that provides a pipeline-based architecture for processing templates with functions and data. The core functionality includes processing templates sequentially, supporting custom template functions, and collecting processed templates. It features a sophisticated template-aware function system that allows functions to have access to the template context.

## Repository Structure
```
/Volumes/Data/workspace/github/pasdam/go-sdk/
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── README.md
├── .build/...
├── .git/...
├── .github/
│   ├── dependabot.yml
│   └── workflows/
├── pkg/
│   ├── collectors/
│   ├── filters/
│   ├── pipeline/        # Core pipeline processing logic
│   ├── templateproviders/
│   ├── templates/       # Template processing utilities
│   └── values/
└── scripts/
```

## Key Components

### Pipeline Package (`pkg/pipeline/`)
- **Pipeline Interface**: Defines the processing contract with `Process(processData map[string]interface{}) error`
- **Template Processing**: Sequential processing of templates using Go's text/template package
- **Function Maps**: Support for custom template functions via `template.FuncMap`
- **Template-Aware Functions**: Special functions that have access to the template context during processing
- **Template Providers**: Interfaces for providing templates to process
- **Collectors**: Interfaces for collecting processed templates

### Core Files
- `pipeline.go`: Main pipeline implementation with processing logic and support for both regular and template-aware functions
- `process_next_template.go`: Logic for processing individual templates, enhanced to support template-aware functions
- `template_provider.go`: Interface for template access (`NextTemplate()`)
- `pipeline_builder.go`: Builder pattern for creating pipeline instances, extended to support template-aware functions
- `template.go`: Template data structure with Reader and Path

### Templates Package (`pkg/templates/`)
- `process_template.go`: Template processing with support for both regular and template-aware functions
- `apply_template.go`: Template execution with both regular and template-aware functions
- **TemplateAwareFuncMap**: Type definition for template-aware function generators

## Template-Aware Functions System

The SDK supports template-aware functions that are functions which can access and interact with the template context. These are defined as:

```go
type TemplateAwareFuncMap map[string]func(*template.Template) any
```

These functions receive the template instance as a parameter, allowing them to:
- Access the current template context
- Perform operations with template-specific information
- Interact with the template processing environment

## Important Interfaces

### TemplateProvider
```go
type TemplateProvider interface {
    NextTemplate() (*Template, error)
}
```

### Collector
```go
type Collector interface {
    Collect(args *Template) error
    OnPipelineCompleted() error
}
```

### Pipeline
```go
type Pipeline interface {
    Process(processData map[string]interface{}) error
}
```

### PipelineBuilder (Extended)
The builder now supports template-aware functions:
```go
type PipelineBuilder interface {
    // ... existing methods ...
    WithTemplateAwareFunctions(functions templates.TemplateAwareFuncMap) *pipelineBuilder
}
```

## Important Patterns & Conventions

### Architecture
- Pipeline-based processing with sequential template iteration
- Builder pattern for pipeline construction (`NewPipelineBuilder()`)
- Dependency injection through interfaces
- Function maps for template customization via `template.FuncMap`
- Template-aware functions via `templates.TemplateAwareFuncMap`

### Error Handling
- EOF errors indicate end of template processing
- Propagate errors from template providers and collectors
- Validation at pipeline construction time

### Testing
- Mock-based testing using testify/assert
- Comprehensive test coverage for pipeline functionality
- Testing of the builder pattern and error conditions

## Extension Points
- New template providers can be implemented by extending the TemplateProvider interface
- Custom collectors can be created by implementing the Collector interface
- Regular template functions can be added via the template.FuncMap mechanism
- Template-aware functions can be added via the templates.TemplateAwareFuncMap mechanism
- Data preprocessors can be added to transform input data

## Version Information
- This is version 2 of the go-sdk (as evident from the go.mod file)
- The import path is `github.com/go-scaffold/go-sdk/v2`