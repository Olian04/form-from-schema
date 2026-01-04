package jsonschema

import (
	"fmt"

	"github.com/Olian04/form-from-schema/lib"
)

// ConvertSchemaToForm converts a JSON Schema to a Form structure
func ConvertSchemaToForm(schema *Schema) (*lib.Form, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	form := &lib.Form{
		Title:       schema.Title,
		Description: schema.Description,
		Method:      "POST", // Default
		Fields:      []lib.Field{},
	}

	// Handle object schemas with properties
	if schema.Properties != nil {
		fields, err := convertPropertiesToFields(schema.Properties, schema.Required)
		if err != nil {
			return nil, err
		}
		form.Fields = fields
	} else {
		// Handle single field schemas
		field, err := convertSchemaToField("", schema)
		if err != nil {
			return nil, err
		}
		if field != nil {
			form.Fields = []lib.Field{*field}
		}
	}

	return form, nil
}

// convertPropertiesToFields converts schema properties to form fields
func convertPropertiesToFields(properties map[string]*Schema, required []string) ([]lib.Field, error) {
	requiredMap := make(map[string]bool)
	for _, req := range required {
		requiredMap[req] = true
	}

	fields := make([]lib.Field, 0, len(properties))
	for name, propSchema := range properties {
		field, err := convertSchemaToField(name, propSchema)
		if err != nil {
			return nil, fmt.Errorf("error converting field %s: %w", name, err)
		}
		if field != nil {
			if requiredMap[name] {
				if field.Validation == nil {
					field.Validation = &lib.Validation{}
				}
				field.Validation.Required = true
			}
			fields = append(fields, *field)
		}
	}

	return fields, nil
}

// convertSchemaToField converts a single schema to a form field
func convertSchemaToField(name string, schema *Schema) (*lib.Field, error) {
	if schema == nil {
		return nil, nil
	}

	field := &lib.Field{
		Name:        name,
		Label:       schema.Title,
		Description: schema.Description,
		Default:     schema.Default,
		ReadOnly:    schema.ReadOnly != nil && *schema.ReadOnly,
		Deprecated:  schema.Deprecated != nil && *schema.Deprecated,
	}

	// Determine field type
	fieldType, err := determineFieldType(schema)
	if err != nil {
		return nil, err
	}
	field.Type = fieldType

	// Handle enum/const - convert to select or radio
	if len(schema.Enum) > 0 {
		field.Options = convertEnumToOptions(schema.Enum)
		if len(field.Options) <= 3 {
			field.Type = lib.FieldTypeRadio
		} else {
			field.Type = lib.FieldTypeSelect
		}
	} else if schema.Const != nil {
		field.Value = schema.Const
		field.Type = lib.FieldTypeHidden
	}

	// Build validation rules
	field.Validation = buildValidation(schema)

	// Handle object type - nested fields
	if fieldType == lib.FieldTypeObject && schema.Properties != nil {
		nestedFields, err := convertPropertiesToFields(schema.Properties, schema.Required)
		if err != nil {
			return nil, err
		}
		field.Fields = nestedFields
	}

	// Handle array type
	if fieldType == lib.FieldTypeArray {
		if schema.Items != nil {
			itemField, err := convertSchemaToField("item", schema.Items)
			if err != nil {
				return nil, err
			}
			if itemField != nil {
				field.Fields = []lib.Field{*itemField}
			}
		}
		if schema.MinItems != nil {
			if field.Validation == nil {
				field.Validation = &lib.Validation{}
			}
			field.Validation.MinItems = schema.MinItems
		}
		if schema.MaxItems != nil {
			if field.Validation == nil {
				field.Validation = &lib.Validation{}
			}
			field.Validation.MaxItems = schema.MaxItems
		}
	}

	// Handle conditional fields (if/then/else)
	if schema.If != nil {
		conditional, err := buildConditionalField(schema)
		if err != nil {
			return nil, err
		}
		if conditional != nil {
			field.Conditional = conditional
		}
	}

	return field, nil
}

// determineFieldType determines the HTML field type from the schema
func determineFieldType(schema *Schema) (lib.FieldType, error) {
	typeStr, typeArray, hasType := schema.GetType()

	if !hasType {
		// If no type specified, try to infer from other properties
		if schema.Properties != nil {
			return lib.FieldTypeObject, nil
		}
		if schema.Items != nil {
			return lib.FieldTypeArray, nil
		}
		return lib.FieldTypeText, nil // Default to text
	}

	// Handle array of types (union types)
	if len(typeArray) > 0 {
		// For union types, prefer the first non-null type
		for _, t := range typeArray {
			if t != "null" {
				return mapJSONTypeToFieldType(t, schema)
			}
		}
		return lib.FieldTypeText, nil
	}

	return mapJSONTypeToFieldType(typeStr, schema)
}

// mapJSONTypeToFieldType maps JSON Schema types to HTML field types
func mapJSONTypeToFieldType(jsonType string, schema *Schema) (lib.FieldType, error) {
	switch jsonType {
	case "string":
		// Check format for more specific types
		switch schema.Format {
		case "email":
			return lib.FieldTypeEmail, nil
		case "uri", "url":
			return lib.FieldTypeURL, nil
		case "date":
			return lib.FieldTypeDate, nil
		case "time":
			return lib.FieldTypeTime, nil
		case "date-time":
			return lib.FieldTypeDateTime, nil
		case "password":
			return lib.FieldTypePassword, nil
		default:
			// Check if it's a long text field (textarea)
			if schema.MaxLength != nil && *schema.MaxLength > 100 {
				return lib.FieldTypeTextarea, nil
			}
			return lib.FieldTypeText, nil
		}
	case "number", "integer":
		return lib.FieldTypeNumber, nil
	case "boolean":
		return lib.FieldTypeCheckbox, nil
	case "array":
		return lib.FieldTypeArray, nil
	case "object":
		return lib.FieldTypeObject, nil
	case "null":
		return lib.FieldTypeText, nil // Default fallback
	default:
		return lib.FieldTypeText, nil
	}
}

// convertEnumToOptions converts enum values to Option structs
func convertEnumToOptions(enum []any) []lib.Option {
	options := make([]lib.Option, 0, len(enum))
	for _, value := range enum {
		options = append(options, lib.Option{
			Label: fmt.Sprintf("%v", value),
			Value: value,
		})
	}
	return options
}

// buildValidation builds validation rules from schema
func buildValidation(schema *Schema) *lib.Validation {
	if schema == nil {
		return nil
	}

	validation := &lib.Validation{}

	// String validations
	if schema.MinLength != nil {
		validation.MinLength = schema.MinLength
	}
	if schema.MaxLength != nil {
		validation.MaxLength = schema.MaxLength
	}
	if schema.Pattern != "" {
		validation.Pattern = schema.Pattern
		validation.PatternError = "Invalid format"
	}

	// Number validations
	if schema.Minimum != nil {
		validation.Min = schema.Minimum
	}
	if schema.Maximum != nil {
		validation.Max = schema.Maximum
	}
	if schema.ExclusiveMinimum != nil {
		validation.Min = schema.ExclusiveMinimum
	}
	if schema.ExclusiveMaximum != nil {
		validation.Max = schema.ExclusiveMaximum
	}
	if schema.MultipleOf != nil {
		validation.Step = schema.MultipleOf
	}

	// Check if validation has any rules
	if validation.MinLength == nil && validation.MaxLength == nil &&
		validation.Min == nil && validation.Max == nil &&
		validation.Pattern == "" && validation.Step == nil {
		return nil
	}

	return validation
}

// buildConditionalField builds conditional field logic from if/then/else
func buildConditionalField(schema *Schema) (*lib.ConditionalField, error) {
	if schema.If == nil {
		return nil, nil
	}

	conditional := &lib.ConditionalField{}

	// Try to extract condition from If schema
	// This is a simplified version - full implementation would need to parse the condition
	if schema.If.Properties != nil {
		// Extract first property as condition field
		for name := range schema.If.Properties {
			conditional.Condition = name
			break
		}
	}

	// Convert Then fields
	if schema.Then != nil {
		thenFields, err := convertPropertiesToFields(schema.Then.Properties, schema.Then.Required)
		if err != nil {
			return nil, err
		}
		conditional.Then = thenFields
	}

	// Convert Else fields
	if schema.Else != nil {
		elseFields, err := convertPropertiesToFields(schema.Else.Properties, schema.Else.Required)
		if err != nil {
			return nil, err
		}
		conditional.Else = elseFields
	}

	if conditional.Condition == "" {
		return nil, nil
	}

	return conditional, nil
}
