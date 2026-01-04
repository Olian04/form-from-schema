package lib

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldType represents the HTML input type for a form field
type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeEmail    FieldType = "email"
	FieldTypePassword FieldType = "password"
	FieldTypeNumber   FieldType = "number"
	FieldTypeTel      FieldType = "tel"
	FieldTypeURL      FieldType = "url"
	FieldTypeDate     FieldType = "date"
	FieldTypeTime     FieldType = "time"
	FieldTypeDateTime FieldType = "datetime-local"
	FieldTypeMonth    FieldType = "month"
	FieldTypeWeek     FieldType = "week"
	FieldTypeTextarea FieldType = "textarea"
	FieldTypeSelect   FieldType = "select"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeRadio    FieldType = "radio"
	FieldTypeFile     FieldType = "file"
	FieldTypeHidden   FieldType = "hidden"
	FieldTypeObject   FieldType = "object" // For nested objects
	FieldTypeArray    FieldType = "array"  // For arrays
)

// Option represents an option for select, radio, or checkbox fields
type Option struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// Validation represents validation rules for a form field
type Validation struct {
	Required     bool     `json:"required,omitempty"`
	MinLength    *int     `json:"minLength,omitempty"`
	MaxLength    *int     `json:"maxLength,omitempty"`
	Min          *float64 `json:"min,omitempty"`
	Max          *float64 `json:"max,omitempty"`
	Pattern      string   `json:"pattern,omitempty"`
	PatternError string   `json:"patternError,omitempty"`
	Step         *float64 `json:"step,omitempty"`
	MinItems     *int     `json:"minItems,omitempty"`
	MaxItems     *int     `json:"maxItems,omitempty"`
}

// ConditionalField represents a conditional field (if/then/else logic)
type ConditionalField struct {
	Condition string  `json:"condition"`      // Field name that triggers this condition
	Value     any     `json:"value"`          // Value that triggers this condition
	Then      []Field `json:"then"`           // Fields to show when condition is met
	Else      []Field `json:"else,omitempty"` // Fields to show when condition is not met
}

// Field represents a single form field
type Field struct {
	Name        string            `json:"name"`
	Type        FieldType         `json:"type"`
	Label       string            `json:"label,omitempty"`
	Description string            `json:"description,omitempty"`
	Placeholder string            `json:"placeholder,omitempty"`
	Default     any               `json:"default,omitempty"`
	Value       any               `json:"value,omitempty"`
	Options     []Option          `json:"options,omitempty"`
	Validation  *Validation       `json:"validation,omitempty"`
	ReadOnly    bool              `json:"readOnly,omitempty"`
	Deprecated  bool              `json:"deprecated,omitempty"`
	Fields      []Field           `json:"fields,omitempty"` // For object/array types
	Conditional *ConditionalField `json:"conditional,omitempty"`
	HelpText    string            `json:"helpText,omitempty"`
}

// Form represents a complete HTML form structure
type Form struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Action      string  `json:"action,omitempty"`
	Method      string  `json:"method,omitempty"`
	Fields      []Field `json:"fields"`
}

// Validate validates the form structure to ensure it's in a valid state
// and can be safely used to generate HTML forms deterministically
func (f *Form) Validate() error {
	if f == nil {
		return fmt.Errorf("form cannot be nil")
	}

	// Validate form has at least one field
	if len(f.Fields) == 0 {
		return fmt.Errorf("form must have at least one field")
	}

	// Validate HTTP method if specified
	if f.Method != "" {
		validMethods := map[string]bool{
			"GET": true, "POST": true, "PUT": true, "PATCH": true,
			"DELETE": true, "HEAD": true, "OPTIONS": true,
		}
		if !validMethods[strings.ToUpper(f.Method)] {
			return fmt.Errorf("invalid HTTP method: %s (must be one of: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS)", f.Method)
		}
	}

	// Track field names to ensure uniqueness
	fieldNames := make(map[string]bool)

	// Validate all top-level fields
	for i, field := range f.Fields {
		if err := f.validateField(&field, fieldNames, fmt.Sprintf("fields[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

// validateField validates a single field and its nested fields recursively
func (f *Form) validateField(field *Field, parentFieldNames map[string]bool, path string) error {
	if field == nil {
		return fmt.Errorf("%s: field cannot be nil", path)
	}

	// Validate field name
	if err := validateFieldName(field.Name, path); err != nil {
		return err
	}

	// Check for duplicate field names at the same level
	if field.Name != "" {
		if parentFieldNames[field.Name] {
			return fmt.Errorf("%s: duplicate field name '%s' at the same level", path, field.Name)
		}
		parentFieldNames[field.Name] = true
	}

	// Validate field type
	if err := validateFieldType(field.Type, path); err != nil {
		return err
	}

	// Validate field type-specific constraints
	if err := f.validateFieldTypeConstraints(field, path); err != nil {
		return err
	}

	// Validate validation rules
	if field.Validation != nil {
		if err := validateValidationRules(field.Validation, field.Type, path); err != nil {
			return err
		}
	}

	// Validate conditional fields
	if field.Conditional != nil {
		if err := f.validateConditionalField(field.Conditional, parentFieldNames, path); err != nil {
			return err
		}
	}

	// Validate nested fields (for objects and arrays)
	if len(field.Fields) > 0 {
		if field.Type != FieldTypeObject && field.Type != FieldTypeArray {
			return fmt.Errorf("%s: fields with nested Fields must have type 'object' or 'array', got '%s'", path, field.Type)
		}

		// Create a new scope for nested field names
		nestedFieldNames := make(map[string]bool)
		for i, nestedField := range field.Fields {
			nestedPath := fmt.Sprintf("%s.fields[%d]", path, i)
			if err := f.validateField(&nestedField, nestedFieldNames, nestedPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateFieldName validates that a field name is valid for HTML forms
func validateFieldName(name string, path string) error {
	if name == "" {
		// Empty names are allowed for top-level single fields, but warn
		return nil
	}

	// HTML form field names must start with a letter or underscore, and contain only
	// letters, digits, underscores, hyphens, and dots
	validNamePattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.-]*$`)
	if !validNamePattern.MatchString(name) {
		return fmt.Errorf("%s: invalid field name '%s' (must start with letter/underscore and contain only letters, digits, underscores, hyphens, and dots)", path, name)
	}

	// Reserved HTML form field names that could cause conflicts
	reservedNames := map[string]bool{
		"submit": true, "reset": true, "button": true,
		"form": true, "fieldset": true, "legend": true,
	}
	if reservedNames[strings.ToLower(name)] {
		return fmt.Errorf("%s: field name '%s' is reserved and cannot be used", path, name)
	}

	return nil
}

// validateFieldType validates that the field type is valid
func validateFieldType(fieldType FieldType, path string) error {
	validTypes := map[FieldType]bool{
		FieldTypeText:     true,
		FieldTypeEmail:    true,
		FieldTypePassword: true,
		FieldTypeNumber:   true,
		FieldTypeTel:      true,
		FieldTypeURL:      true,
		FieldTypeDate:     true,
		FieldTypeTime:     true,
		FieldTypeDateTime: true,
		FieldTypeMonth:    true,
		FieldTypeWeek:     true,
		FieldTypeTextarea: true,
		FieldTypeSelect:   true,
		FieldTypeCheckbox: true,
		FieldTypeRadio:    true,
		FieldTypeFile:     true,
		FieldTypeHidden:   true,
		FieldTypeObject:   true,
		FieldTypeArray:    true,
	}

	if !validTypes[fieldType] {
		return fmt.Errorf("%s: invalid field type '%s'", path, fieldType)
	}

	return nil
}

// validateFieldTypeConstraints validates type-specific constraints
func (f *Form) validateFieldTypeConstraints(field *Field, path string) error {
	// Select and Radio fields must have options
	if field.Type == FieldTypeSelect || field.Type == FieldTypeRadio {
		if len(field.Options) == 0 {
			return fmt.Errorf("%s: field type '%s' requires at least one option", path, field.Type)
		}

		// Validate option values are unique
		optionValues := make(map[string]bool)
		for i, option := range field.Options {
			optionValue := fmt.Sprintf("%v", option.Value)
			if optionValues[optionValue] {
				return fmt.Errorf("%s: duplicate option value '%s' at options[%d]", path, optionValue, i)
			}
			optionValues[optionValue] = true
		}
	}

	// Checkbox fields typically shouldn't have options (they're boolean)
	// But we'll allow it for multi-select scenarios
	if field.Type == FieldTypeCheckbox && len(field.Options) > 0 {
		// Validate option values are unique
		optionValues := make(map[string]bool)
		for i, option := range field.Options {
			optionValue := fmt.Sprintf("%v", option.Value)
			if optionValues[optionValue] {
				return fmt.Errorf("%s: duplicate option value '%s' at options[%d]", path, optionValue, i)
			}
			optionValues[optionValue] = true
		}
	}

	// Text, Email, Password, URL, Tel fields shouldn't have options
	if field.Type == FieldTypeText || field.Type == FieldTypeEmail ||
		field.Type == FieldTypePassword || field.Type == FieldTypeURL ||
		field.Type == FieldTypeTel || field.Type == FieldTypeTextarea {
		if len(field.Options) > 0 {
			return fmt.Errorf("%s: field type '%s' cannot have options", path, field.Type)
		}
	}

	// Number fields shouldn't have options
	if field.Type == FieldTypeNumber {
		if len(field.Options) > 0 {
			return fmt.Errorf("%s: field type '%s' cannot have options", path, field.Type)
		}
	}

	// Date/Time fields shouldn't have options
	if field.Type == FieldTypeDate || field.Type == FieldTypeTime ||
		field.Type == FieldTypeDateTime || field.Type == FieldTypeMonth ||
		field.Type == FieldTypeWeek {
		if len(field.Options) > 0 {
			return fmt.Errorf("%s: field type '%s' cannot have options", path, field.Type)
		}
	}

	return nil
}

// validateValidationRules validates that validation rules are consistent
func validateValidationRules(validation *Validation, fieldType FieldType, path string) error {
	if validation == nil {
		return nil
	}

	// Validate string length constraints
	if validation.MinLength != nil {
		if *validation.MinLength < 0 {
			return fmt.Errorf("%s: validation.minLength cannot be negative", path)
		}
	}
	if validation.MaxLength != nil {
		if *validation.MaxLength < 0 {
			return fmt.Errorf("%s: validation.maxLength cannot be negative", path)
		}
	}
	if validation.MinLength != nil && validation.MaxLength != nil {
		if *validation.MinLength > *validation.MaxLength {
			return fmt.Errorf("%s: validation.minLength (%d) cannot be greater than maxLength (%d)", path, *validation.MinLength, *validation.MaxLength)
		}
	}

	// Validate number range constraints
	if validation.Min != nil && validation.Max != nil {
		if *validation.Min > *validation.Max {
			return fmt.Errorf("%s: validation.min (%v) cannot be greater than max (%v)", path, *validation.Min, *validation.Max)
		}
	}

	// Validate array item constraints
	if validation.MinItems != nil {
		if *validation.MinItems < 0 {
			return fmt.Errorf("%s: validation.minItems cannot be negative", path)
		}
	}
	if validation.MaxItems != nil {
		if *validation.MaxItems < 0 {
			return fmt.Errorf("%s: validation.maxItems cannot be negative", path)
		}
	}
	if validation.MinItems != nil && validation.MaxItems != nil {
		if *validation.MinItems > *validation.MaxItems {
			return fmt.Errorf("%s: validation.minItems (%d) cannot be greater than maxItems (%d)", path, *validation.MinItems, *validation.MaxItems)
		}
	}

	// Validate step is positive
	if validation.Step != nil {
		if *validation.Step <= 0 {
			return fmt.Errorf("%s: validation.step must be positive, got %v", path, *validation.Step)
		}
	}

	// Validate that validation rules match field type
	if fieldType == FieldTypeText || fieldType == FieldTypeEmail ||
		fieldType == FieldTypePassword || fieldType == FieldTypeURL ||
		fieldType == FieldTypeTel || fieldType == FieldTypeTextarea {
		// String validations
		if validation.Min != nil || validation.Max != nil || validation.Step != nil {
			return fmt.Errorf("%s: validation rules min/max/step are not applicable for field type '%s'", path, fieldType)
		}
		if validation.MinItems != nil || validation.MaxItems != nil {
			return fmt.Errorf("%s: validation rules minItems/maxItems are not applicable for field type '%s'", path, fieldType)
		}
	}

	if fieldType == FieldTypeNumber {
		// Number validations
		if validation.MinLength != nil || validation.MaxLength != nil {
			return fmt.Errorf("%s: validation rules minLength/maxLength are not applicable for field type '%s'", path, fieldType)
		}
		if validation.MinItems != nil || validation.MaxItems != nil {
			return fmt.Errorf("%s: validation rules minItems/maxItems are not applicable for field type '%s'", path, fieldType)
		}
	}

	if fieldType == FieldTypeArray {
		// Array validations
		if validation.MinLength != nil || validation.MaxLength != nil {
			return fmt.Errorf("%s: validation rules minLength/maxLength are not applicable for field type '%s'", path, fieldType)
		}
		if validation.Min != nil || validation.Max != nil || validation.Step != nil {
			return fmt.Errorf("%s: validation rules min/max/step are not applicable for field type '%s'", path, fieldType)
		}
	}

	if fieldType == FieldTypeCheckbox || fieldType == FieldTypeRadio || fieldType == FieldTypeSelect {
		// These types typically don't use numeric or length validations
		if validation.MinLength != nil || validation.MaxLength != nil {
			return fmt.Errorf("%s: validation rules minLength/maxLength are not applicable for field type '%s'", path, fieldType)
		}
		if validation.Min != nil || validation.Max != nil || validation.Step != nil {
			return fmt.Errorf("%s: validation rules min/max/step are not applicable for field type '%s'", path, fieldType)
		}
		if validation.MinItems != nil || validation.MaxItems != nil {
			return fmt.Errorf("%s: validation rules minItems/maxItems are not applicable for field type '%s'", path, fieldType)
		}
	}

	return nil
}

// validateConditionalField validates conditional field logic
func (f *Form) validateConditionalField(conditional *ConditionalField, parentFieldNames map[string]bool, path string) error {
	if conditional == nil {
		return nil
	}

	// Validate condition field name exists
	if conditional.Condition == "" {
		return fmt.Errorf("%s: conditional field must specify a condition field name", path)
	}

	// Check that the condition field exists in the parent scope
	if !parentFieldNames[conditional.Condition] {
		return fmt.Errorf("%s: conditional field references non-existent field '%s'", path, conditional.Condition)
	}

	// Validate Then fields
	if len(conditional.Then) > 0 {
		thenFieldNames := make(map[string]bool)
		for i, field := range conditional.Then {
			thenPath := fmt.Sprintf("%s.conditional.then[%d]", path, i)
			if err := f.validateField(&field, thenFieldNames, thenPath); err != nil {
				return err
			}
		}
	}

	// Validate Else fields
	if len(conditional.Else) > 0 {
		elseFieldNames := make(map[string]bool)
		for i, field := range conditional.Else {
			elsePath := fmt.Sprintf("%s.conditional.else[%d]", path, i)
			if err := f.validateField(&field, elseFieldNames, elsePath); err != nil {
				return err
			}
		}
	}

	return nil
}
