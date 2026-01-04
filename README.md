# Form From Schema

This project provides a generic schema-to-form converter with a focus on JSON Schema to HTML form conversion.

## Features

- **JSON Schema Support**: Full support for JSON Schema Draft 2020-12 specification
- **Type Mapping**: Automatic mapping of JSON Schema types to appropriate HTML input types
- **Validation Rules**: Converts JSON Schema validation rules to HTML form validation attributes
- **Nested Structures**: Support for objects and arrays with nested fields
- **Conditional Fields**: Support for conditional field display using `if/then/else` logic
- **Form Validation**: Comprehensive validation to ensure forms are valid and deterministic
- **Type Safety**: Strongly typed Go structures for form definitions
- **HTML Generation**: Uses [templ](https://templ.guide/) for type-safe HTML generation

## Installation

```bash
go get github.com/Olian04/form-from-schema
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/Olian04/form-from-schema"
)

func main() {
    // JSON Schema definition
    schema := []byte(`{
        "type": "object",
        "title": "User Registration",
        "description": "Register a new user account",
        "properties": {
            "username": {
                "type": "string",
                "title": "Username",
                "minLength": 3,
                "maxLength": 20
            },
            "email": {
                "type": "string",
                "format": "email",
                "title": "Email Address"
            },
            "age": {
                "type": "integer",
                "title": "Age",
                "minimum": 18,
                "maximum": 120
            }
        },
        "required": ["username", "email"]
    }`)

    // Convert JSON Schema to Form
    form, err := formfromschema.FromJsonSchema(schema)
    if err != nil {
        panic(err)
    }

    form.Action = "/api/resource/update"
    form.Method = "POST"

    // Validate the form
    if err := form.Validate(); err != nil {
        panic(err)
    }

    // Generate HTML
    ctx := context.Background()
    if err := formfromschema.ToHtml(ctx, form, os.Stdout); err != nil {
        panic(err)
    }
}
```

## Usage

### Converting JSON Schema to Form

```go
import "github.com/Olian04/form-from-schema"

schema := []byte(`{
    "type": "object",
    "properties": {
        "name": {"type": "string"}
    }
}`)

form, err := formfromschema.FromJsonSchema(schema)
if err != nil {
    // handle error
}
```

### Validating Forms

Before generating HTML, it's recommended to validate the form:

```go
if err := form.Validate(); err != nil {
    // Form is invalid - handle error
    log.Fatal(err)
}
```

The validation ensures:

- Field names are unique and valid HTML form field names
- Field types are valid
- Validation rules are consistent (e.g., min ≤ max)
- Conditional fields reference valid fields
- No configuration conflicts

### Generating HTML

```go
import (
    "context"
    "os"
    "github.com/Olian04/form-from-schema"
)

ctx := context.Background()
err := formfromschema.ToHtml(ctx, form, os.Stdout)
```

You can write to any `io.Writer`:

```go
var buf bytes.Buffer
err := formfromschema.ToHtml(ctx, form, &buf)
html := buf.String()
```

## Supported Field Types

The library automatically maps JSON Schema types to HTML input types:

| JSON Schema Type | HTML Input Type | Notes |
|-----------------|----------------|-------|
| `string` | `text` | Default for string types |
| `string` (format: `email`) | `email` | Email input |
| `string` (format: `url`) | `url` | URL input |
| `string` (format: `date`) | `date` | Date picker |
| `string` (format: `date-time`) | `datetime-local` | Date and time picker |
| `string` (format: `time`) | `time` | Time picker |
| `string` (format: `password`) | `password` | Password input |
| `string` (maxLength > 100) | `textarea` | Long text fields |
| `number`, `integer` | `number` | Number input |
| `boolean` | `checkbox` | Checkbox |
| `array` | `array` | Array with nested item fields |
| `object` | `object` | Object with nested fields |

### Enum Handling

- **2-3 enum values**: Converted to radio buttons
- **4+ enum values**: Converted to select dropdown
- **const value**: Converted to hidden input

## JSON Schema Features

### Supported JSON Schema Features

- ✅ Core vocabulary (`$schema`, `$id`, `$ref`, `$defs`, etc.)
- ✅ Applicator vocabulary (`allOf`, `anyOf`, `oneOf`, `if/then/else`, etc.)
- ✅ Validation vocabulary (all validation keywords)
- ✅ Meta-data vocabulary (`title`, `description`, `default`, etc.)
- ✅ Nested objects and arrays
- ✅ Conditional fields (`if/then/else`)
- ✅ Enum and const values
- ✅ Format annotations (email, url, date, etc.)

### Example: Complex Schema

```json
{
  "type": "object",
  "title": "User Profile",
  "properties": {
    "username": {
      "type": "string",
      "title": "Username",
      "minLength": 3,
      "maxLength": 20,
      "pattern": "^[a-zA-Z0-9_]+$"
    },
    "email": {
      "type": "string",
      "format": "email",
      "title": "Email"
    },
    "country": {
      "type": "string",
      "enum": ["USA", "Canada", "UK"],
      "title": "Country"
    },
    "address": {
      "type": "object",
      "title": "Address",
      "properties": {
        "street": {"type": "string"},
        "city": {"type": "string"},
        "zip": {"type": "string"}
      }
    },
    "tags": {
      "type": "array",
      "items": {"type": "string"},
      "minItems": 1,
      "maxItems": 10
    }
  },
  "required": ["username", "email"]
}
```

## Form Structure

The `lib.Form` struct represents a form:

```go
type Form struct {
    Title       string  // Form title
    Description string  // Form description
    Action      string  // Form action URL
    Method      string  // HTTP method (GET, POST, etc.)
    Fields      []Field // Form fields
}
```

Each field supports:

- **Name**: HTML field name (must be unique)
- **Type**: Field type (text, email, number, etc.)
- **Label**: Display label
- **Description**: Help text
- **Placeholder**: Placeholder text
- **Default**: Default value
- **Options**: Options for select/radio fields
- **Validation**: Validation rules
- **Fields**: Nested fields (for objects/arrays)
- **Conditional**: Conditional field logic

## Validation

The `Form.Validate()` method performs comprehensive validation:

- ✅ Ensures field names are unique
- ✅ Validates HTML field name format
- ✅ Checks for reserved HTML names
- ✅ Validates field types
- ✅ Ensures validation rules are consistent
- ✅ Validates conditional field references
- ✅ Validates nested field structures

Example:

```go
if err := form.Validate(); err != nil {
    fmt.Printf("Validation error: %v\n", err)
    // Output: Validation error: fields[1]: duplicate field name 'username' at the same level
}
```

## Architecture

The project is organized into several packages:

```
lib/
├── form.go              # Core Form and Field types
├── schemas/
│   └── jsonschema/      # JSON Schema parsing and conversion
│       ├── schema.go   # JSON Schema types
│       └── convert.go   # Schema to Form conversion
└── targets/
    └── html/            # HTML form generation
        ├── form.templ   # Form template
        ├── field.templ  # Field template
        └── convert.go   # Form to HTML conversion
```

### Design Principles

1. **Separation of Concerns**: Schema parsing, form representation, and HTML generation are separate
2. **Extensibility**: Easy to add new schema formats or output targets
3. **Type Safety**: Strong typing throughout the codebase
4. **Validation**: Forms are validated before HTML generation

## Development

### Prerequisites

- Go 1.25.5 or later
- [templ](https://templ.guide/) for HTML template generation

### Building

```bash
# Generate templ files
make generate

# Run tests
make test

# Format code
make format

# Run linter
make lint

# Run everything
make all
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./lib/targets/html -v
```

## Contributing

Contributions are welcome! Please ensure:

1. All tests pass (`make test`)
2. Code is formatted (`make format`)
3. Linter passes (`make lint`)
4. New features include tests
5. Documentation is updated

## License

See [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Support for additional schema formats (OpenAPI, GraphQL, etc.)
- [ ] Additional output targets (HTMX, React components, Vue components, etc.)
- [ ] Custom field type mappings
- [ ] Form builder API
- [ ] Client-side validation code generation
