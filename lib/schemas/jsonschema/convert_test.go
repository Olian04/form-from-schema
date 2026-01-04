package jsonschema

import (
	"encoding/json"
	"testing"

	"github.com/Olian04/form-from-schema/lib"
)

func TestConvertSchemaToForm(t *testing.T) {
	tests := []struct {
		name    string
		schema  *Schema
		wantErr bool
		check   func(*lib.Form) bool
	}{
		{
			name:    "nil schema",
			schema:  nil,
			wantErr: true,
			check:   nil,
		},
		{
			name: "simple string field",
			schema: &Schema{
				Type:  json.RawMessage(`"string"`),
				Title: "Name",
			},
			wantErr: false,
			check: func(f *lib.Form) bool {
				return len(f.Fields) == 1 && f.Fields[0].Type == lib.FieldTypeText
			},
		},
		{
			name: "object with properties",
			schema: &Schema{
				Type: json.RawMessage(`"object"`),
				Properties: map[string]*Schema{
					"name": {
						Type:  json.RawMessage(`"string"`),
						Title: "Name",
					},
					"age": {
						Type:  json.RawMessage(`"integer"`),
						Title: "Age",
					},
				},
				Required: []string{"name"},
			},
			wantErr: false,
			check: func(f *lib.Form) bool {
				if len(f.Fields) != 2 {
					return false
				}
				// Check that name field is required
				for _, field := range f.Fields {
					if field.Name == "name" {
						return field.Validation != nil && field.Validation.Required
					}
				}
				return false
			},
		},
		{
			name: "form with title and description",
			schema: &Schema{
				Title:       "User Form",
				Description: "A form for user information",
				Type:        json.RawMessage(`"object"`),
				Properties: map[string]*Schema{
					"email": {
						Type: json.RawMessage(`"string"`),
					},
				},
			},
			wantErr: false,
			check: func(f *lib.Form) bool {
				return f.Title == "User Form" && f.Description == "A form for user information"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertSchemaToForm(tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertSchemaToForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(got) {
				t.Errorf("ConvertSchemaToForm() check failed")
			}
		})
	}
}

func TestConvertSchemaToForm_FieldTypes(t *testing.T) {
	tests := []struct {
		name       string
		schema     *Schema
		wantType   lib.FieldType
		fieldIndex int
	}{
		{
			name: "email format",
			schema: &Schema{
				Type:   json.RawMessage(`"string"`),
				Format: "email",
			},
			wantType:   lib.FieldTypeEmail,
			fieldIndex: 0,
		},
		{
			name: "url format",
			schema: &Schema{
				Type:   json.RawMessage(`"string"`),
				Format: "url",
			},
			wantType:   lib.FieldTypeURL,
			fieldIndex: 0,
		},
		{
			name: "date format",
			schema: &Schema{
				Type:   json.RawMessage(`"string"`),
				Format: "date",
			},
			wantType:   lib.FieldTypeDate,
			fieldIndex: 0,
		},
		{
			name: "date-time format",
			schema: &Schema{
				Type:   json.RawMessage(`"string"`),
				Format: "date-time",
			},
			wantType:   lib.FieldTypeDateTime,
			fieldIndex: 0,
		},
		{
			name: "password format",
			schema: &Schema{
				Type:   json.RawMessage(`"string"`),
				Format: "password",
			},
			wantType:   lib.FieldTypePassword,
			fieldIndex: 0,
		},
		{
			name: "number type",
			schema: &Schema{
				Type: json.RawMessage(`"number"`),
			},
			wantType:   lib.FieldTypeNumber,
			fieldIndex: 0,
		},
		{
			name: "integer type",
			schema: &Schema{
				Type: json.RawMessage(`"integer"`),
			},
			wantType:   lib.FieldTypeNumber,
			fieldIndex: 0,
		},
		{
			name: "boolean type",
			schema: &Schema{
				Type: json.RawMessage(`"boolean"`),
			},
			wantType:   lib.FieldTypeCheckbox,
			fieldIndex: 0,
		},
		{
			name: "array type",
			schema: &Schema{
				Type: json.RawMessage(`"array"`),
			},
			wantType:   lib.FieldTypeArray,
			fieldIndex: 0,
		},
		{
			name: "object type",
			schema: &Schema{
				Type: json.RawMessage(`"object"`),
			},
			wantType:   lib.FieldTypeObject,
			fieldIndex: 0,
		},
		{
			name: "long text becomes textarea",
			schema: &Schema{
				Type:      json.RawMessage(`"string"`),
				MaxLength: intPtr(200),
			},
			wantType:   lib.FieldTypeTextarea,
			fieldIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form, err := ConvertSchemaToForm(tt.schema)
			if err != nil {
				t.Fatalf("ConvertSchemaToForm() error = %v", err)
			}
			if len(form.Fields) <= tt.fieldIndex {
				t.Fatalf("ConvertSchemaToForm() returned form with %d fields, want at least %d", len(form.Fields), tt.fieldIndex+1)
			}
			if form.Fields[tt.fieldIndex].Type != tt.wantType {
				t.Errorf("ConvertSchemaToForm() field type = %v, want %v", form.Fields[tt.fieldIndex].Type, tt.wantType)
			}
		})
	}
}

func TestConvertSchemaToForm_Enum(t *testing.T) {
	tests := []struct {
		name        string
		schema      *Schema
		wantType    lib.FieldType
		wantOptions int
	}{
		{
			name: "enum with 2 options becomes radio",
			schema: &Schema{
				Type: json.RawMessage(`"string"`),
				Enum: []any{"option1", "option2"},
			},
			wantType:    lib.FieldTypeRadio,
			wantOptions: 2,
		},
		{
			name: "enum with 3 options becomes radio",
			schema: &Schema{
				Type: json.RawMessage(`"string"`),
				Enum: []any{"option1", "option2", "option3"},
			},
			wantType:    lib.FieldTypeRadio,
			wantOptions: 3,
		},
		{
			name: "enum with 4 options becomes select",
			schema: &Schema{
				Type: json.RawMessage(`"string"`),
				Enum: []any{"option1", "option2", "option3", "option4"},
			},
			wantType:    lib.FieldTypeSelect,
			wantOptions: 4,
		},
		{
			name: "enum with mixed types",
			schema: &Schema{
				Type: json.RawMessage(`"string"`),
				Enum: []any{"string", 42, true},
			},
			wantType:    lib.FieldTypeRadio,
			wantOptions: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form, err := ConvertSchemaToForm(tt.schema)
			if err != nil {
				t.Fatalf("ConvertSchemaToForm() error = %v", err)
			}
			if len(form.Fields) == 0 {
				t.Fatalf("ConvertSchemaToForm() returned form with no fields")
			}
			field := form.Fields[0]
			if field.Type != tt.wantType {
				t.Errorf("ConvertSchemaToForm() field type = %v, want %v", field.Type, tt.wantType)
			}
			if len(field.Options) != tt.wantOptions {
				t.Errorf("ConvertSchemaToForm() field options count = %v, want %v", len(field.Options), tt.wantOptions)
			}
		})
	}
}

func TestConvertSchemaToForm_Const(t *testing.T) {
	schema := &Schema{
		Type:  json.RawMessage(`"string"`),
		Const: "fixed-value",
	}

	form, err := ConvertSchemaToForm(schema)
	if err != nil {
		t.Fatalf("ConvertSchemaToForm() error = %v", err)
	}

	if len(form.Fields) == 0 {
		t.Fatalf("ConvertSchemaToForm() returned form with no fields")
	}

	field := form.Fields[0]
	if field.Type != lib.FieldTypeHidden {
		t.Errorf("ConvertSchemaToForm() field type = %v, want %v", field.Type, lib.FieldTypeHidden)
	}
	if field.Value != "fixed-value" {
		t.Errorf("ConvertSchemaToForm() field value = %v, want fixed-value", field.Value)
	}
}

func TestConvertSchemaToForm_Validation(t *testing.T) {
	tests := []struct {
		name   string
		schema *Schema
		check  func(*lib.Validation) bool
	}{
		{
			name: "string min/max length",
			schema: &Schema{
				Type:      json.RawMessage(`"string"`),
				MinLength: intPtr(5),
				MaxLength: intPtr(10),
			},
			check: func(v *lib.Validation) bool {
				return v != nil && v.MinLength != nil && *v.MinLength == 5 &&
					v.MaxLength != nil && *v.MaxLength == 10
			},
		},
		{
			name: "number min/max",
			schema: &Schema{
				Type:    json.RawMessage(`"number"`),
				Minimum: floatPtr(0),
				Maximum: floatPtr(100),
			},
			check: func(v *lib.Validation) bool {
				return v != nil && v.Min != nil && *v.Min == 0 &&
					v.Max != nil && *v.Max == 100
			},
		},
		{
			name: "pattern validation",
			schema: &Schema{
				Type:    json.RawMessage(`"string"`),
				Pattern: "^[A-Z]+$",
			},
			check: func(v *lib.Validation) bool {
				return v != nil && v.Pattern == "^[A-Z]+$" && v.PatternError != ""
			},
		},
		{
			name: "multipleOf step",
			schema: &Schema{
				Type:       json.RawMessage(`"number"`),
				MultipleOf: floatPtr(0.5),
			},
			check: func(v *lib.Validation) bool {
				return v != nil && v.Step != nil && *v.Step == 0.5
			},
		},
		{
			name: "exclusive min/max",
			schema: &Schema{
				Type:             json.RawMessage(`"number"`),
				ExclusiveMinimum: floatPtr(0),
				ExclusiveMaximum: floatPtr(100),
			},
			check: func(v *lib.Validation) bool {
				return v != nil && v.Min != nil && *v.Min == 0 &&
					v.Max != nil && *v.Max == 100
			},
		},
		{
			name: "array min/max items",
			schema: &Schema{
				Type:     json.RawMessage(`"array"`),
				MinItems: intPtr(1),
				MaxItems: intPtr(10),
			},
			check: func(v *lib.Validation) bool {
				return v != nil && v.MinItems != nil && *v.MinItems == 1 &&
					v.MaxItems != nil && *v.MaxItems == 10
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form, err := ConvertSchemaToForm(tt.schema)
			if err != nil {
				t.Fatalf("ConvertSchemaToForm() error = %v", err)
			}
			if len(form.Fields) == 0 {
				t.Fatalf("ConvertSchemaToForm() returned form with no fields")
			}
			if !tt.check(form.Fields[0].Validation) {
				t.Errorf("ConvertSchemaToForm() validation check failed")
			}
		})
	}
}

func TestConvertSchemaToForm_NestedObjects(t *testing.T) {
	schema := &Schema{
		Type: json.RawMessage(`"object"`),
		Properties: map[string]*Schema{
			"user": {
				Type: json.RawMessage(`"object"`),
				Properties: map[string]*Schema{
					"name": {
						Type:  json.RawMessage(`"string"`),
						Title: "Name",
					},
					"email": {
						Type:   json.RawMessage(`"string"`),
						Format: "email",
					},
				},
			},
		},
	}

	form, err := ConvertSchemaToForm(schema)
	if err != nil {
		t.Fatalf("ConvertSchemaToForm() error = %v", err)
	}

	if len(form.Fields) != 1 {
		t.Fatalf("ConvertSchemaToForm() returned form with %d fields, want 1", len(form.Fields))
	}

	userField := form.Fields[0]
	if userField.Type != lib.FieldTypeObject {
		t.Errorf("ConvertSchemaToForm() user field type = %v, want %v", userField.Type, lib.FieldTypeObject)
	}

	if len(userField.Fields) != 2 {
		t.Fatalf("ConvertSchemaToForm() user field has %d nested fields, want 2", len(userField.Fields))
	}
}

func TestConvertSchemaToForm_Arrays(t *testing.T) {
	schema := &Schema{
		Type: json.RawMessage(`"array"`),
		Items: &Schema{
			Type:  json.RawMessage(`"string"`),
			Title: "Item",
		},
		MinItems: intPtr(1),
		MaxItems: intPtr(5),
	}

	form, err := ConvertSchemaToForm(schema)
	if err != nil {
		t.Fatalf("ConvertSchemaToForm() error = %v", err)
	}

	if len(form.Fields) == 0 {
		t.Fatalf("ConvertSchemaToForm() returned form with no fields")
	}

	arrayField := form.Fields[0]
	if arrayField.Type != lib.FieldTypeArray {
		t.Errorf("ConvertSchemaToForm() array field type = %v, want %v", arrayField.Type, lib.FieldTypeArray)
	}

	if len(arrayField.Fields) != 1 {
		t.Fatalf("ConvertSchemaToForm() array field has %d item fields, want 1", len(arrayField.Fields))
	}

	if arrayField.Validation == nil {
		t.Fatalf("ConvertSchemaToForm() array field validation is nil")
	}

	if arrayField.Validation.MinItems == nil || *arrayField.Validation.MinItems != 1 {
		t.Errorf("ConvertSchemaToForm() array field minItems = %v, want 1", arrayField.Validation.MinItems)
	}

	if arrayField.Validation.MaxItems == nil || *arrayField.Validation.MaxItems != 5 {
		t.Errorf("ConvertSchemaToForm() array field maxItems = %v, want 5", arrayField.Validation.MaxItems)
	}
}

func TestConvertSchemaToForm_ReadOnlyAndDeprecated(t *testing.T) {
	readOnly := true
	deprecated := true

	schema := &Schema{
		Type:       json.RawMessage(`"string"`),
		ReadOnly:   &readOnly,
		Deprecated: &deprecated,
	}

	form, err := ConvertSchemaToForm(schema)
	if err != nil {
		t.Fatalf("ConvertSchemaToForm() error = %v", err)
	}

	if len(form.Fields) == 0 {
		t.Fatalf("ConvertSchemaToForm() returned form with no fields")
	}

	field := form.Fields[0]
	if !field.ReadOnly {
		t.Errorf("ConvertSchemaToForm() field ReadOnly = %v, want true", field.ReadOnly)
	}
	if !field.Deprecated {
		t.Errorf("ConvertSchemaToForm() field Deprecated = %v, want true", field.Deprecated)
	}
}

func TestConvertSchemaToForm_DefaultValue(t *testing.T) {
	schema := &Schema{
		Type:    json.RawMessage(`"string"`),
		Default: "default-value",
	}

	form, err := ConvertSchemaToForm(schema)
	if err != nil {
		t.Fatalf("ConvertSchemaToForm() error = %v", err)
	}

	if len(form.Fields) == 0 {
		t.Fatalf("ConvertSchemaToForm() returned form with no fields")
	}

	field := form.Fields[0]
	if field.Default != "default-value" {
		t.Errorf("ConvertSchemaToForm() field Default = %v, want default-value", field.Default)
	}
}

func TestConvertSchemaToForm_UnionTypes(t *testing.T) {
	schema := &Schema{
		Type: json.RawMessage(`["string", "null"]`),
	}

	form, err := ConvertSchemaToForm(schema)
	if err != nil {
		t.Fatalf("ConvertSchemaToForm() error = %v", err)
	}

	if len(form.Fields) == 0 {
		t.Fatalf("ConvertSchemaToForm() returned form with no fields")
	}

	field := form.Fields[0]
	if field.Type != lib.FieldTypeText {
		t.Errorf("ConvertSchemaToForm() union type field type = %v, want %v", field.Type, lib.FieldTypeText)
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
