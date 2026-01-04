package lib

import (
	"strings"
	"testing"
)

func TestForm_Validate(t *testing.T) {
	tests := []struct {
		name    string
		form    *Form
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil form",
			form:    nil,
			wantErr: true,
			errMsg:  "form cannot be nil",
		},
		{
			name: "empty fields",
			form: &Form{
				Fields: []Field{},
			},
			wantErr: true,
			errMsg:  "form must have at least one field",
		},
		{
			name: "valid form with single field",
			form: &Form{
				Fields: []Field{
					{
						Name: "username",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid form with multiple fields",
			form: &Form{
				Fields: []Field{
					{
						Name: "username",
						Type: FieldTypeText,
					},
					{
						Name: "email",
						Type: FieldTypeEmail,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid HTTP method",
			form: &Form{
				Method: "INVALID",
				Fields: []Field{
					{
						Name: "username",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid HTTP method",
		},
		{
			name: "valid HTTP method POST",
			form: &Form{
				Method: "POST",
				Fields: []Field{
					{
						Name: "username",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid HTTP method GET",
			form: &Form{
				Method: "GET",
				Fields: []Field{
					{
						Name: "username",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() == "" || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Form.Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestForm_Validate_FieldNames(t *testing.T) {
	tests := []struct {
		name    string
		form    *Form
		wantErr bool
		errMsg  string
	}{
		{
			name: "duplicate field names",
			form: &Form{
				Fields: []Field{
					{
						Name: "username",
						Type: FieldTypeText,
					},
					{
						Name: "username",
						Type: FieldTypeEmail,
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate field name",
		},
		{
			name: "invalid field name starting with number",
			form: &Form{
				Fields: []Field{
					{
						Name: "123field",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid field name",
		},
		{
			name: "invalid field name with special characters",
			form: &Form{
				Fields: []Field{
					{
						Name: "field@name",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid field name",
		},
		{
			name: "reserved field name 'submit'",
			form: &Form{
				Fields: []Field{
					{
						Name: "submit",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: true,
			errMsg:  "reserved",
		},
		{
			name: "reserved field name 'reset'",
			form: &Form{
				Fields: []Field{
					{
						Name: "reset",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: true,
			errMsg:  "reserved",
		},
		{
			name: "valid field name with underscore",
			form: &Form{
				Fields: []Field{
					{
						Name: "_private_field",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid field name with hyphen",
			form: &Form{
				Fields: []Field{
					{
						Name: "field-name",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid field name with dot",
			form: &Form{
				Fields: []Field{
					{
						Name: "field.name",
						Type: FieldTypeText,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Form.Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestForm_Validate_FieldTypes(t *testing.T) {
	tests := []struct {
		name    string
		form    *Form
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid field type",
			form: &Form{
				Fields: []Field{
					{
						Name: "field",
						Type: FieldType("invalid"),
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid field type",
		},
		{
			name: "select field without options",
			form: &Form{
				Fields: []Field{
					{
						Name:    "choice",
						Type:    FieldTypeSelect,
						Options: []Option{},
					},
				},
			},
			wantErr: true,
			errMsg:  "requires at least one option",
		},
		{
			name: "radio field without options",
			form: &Form{
				Fields: []Field{
					{
						Name:    "choice",
						Type:    FieldTypeRadio,
						Options: []Option{},
					},
				},
			},
			wantErr: true,
			errMsg:  "requires at least one option",
		},
		{
			name: "select field with duplicate option values",
			form: &Form{
				Fields: []Field{
					{
						Name: "choice",
						Type: FieldTypeSelect,
						Options: []Option{
							{Label: "Option 1", Value: "opt1"},
							{Label: "Option 2", Value: "opt1"},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate option value",
		},
		{
			name: "text field with options",
			form: &Form{
				Fields: []Field{
					{
						Name: "text",
						Type: FieldTypeText,
						Options: []Option{
							{Label: "Option 1", Value: "opt1"},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot have options",
		},
		{
			name: "number field with options",
			form: &Form{
				Fields: []Field{
					{
						Name: "number",
						Type: FieldTypeNumber,
						Options: []Option{
							{Label: "Option 1", Value: "opt1"},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot have options",
		},
		{
			name: "valid select field with options",
			form: &Form{
				Fields: []Field{
					{
						Name: "choice",
						Type: FieldTypeSelect,
						Options: []Option{
							{Label: "Option 1", Value: "opt1"},
							{Label: "Option 2", Value: "opt2"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "field with nested fields but wrong type",
			form: &Form{
				Fields: []Field{
					{
						Name: "field",
						Type: FieldTypeText,
						Fields: []Field{
							{Name: "nested", Type: FieldTypeText},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "must have type 'object' or 'array'",
		},
		{
			name: "valid object field with nested fields",
			form: &Form{
				Fields: []Field{
					{
						Name: "user",
						Type: FieldTypeObject,
						Fields: []Field{
							{Name: "name", Type: FieldTypeText},
							{Name: "email", Type: FieldTypeEmail},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Form.Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestForm_Validate_ValidationRules(t *testing.T) {
	tests := []struct {
		name    string
		form    *Form
		wantErr bool
		errMsg  string
	}{
		{
			name: "minLength greater than maxLength",
			form: &Form{
				Fields: []Field{
					{
						Name: "text",
						Type: FieldTypeText,
						Validation: &Validation{
							MinLength: intPtr(10),
							MaxLength: intPtr(5),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot be greater than maxLength",
		},
		{
			name: "minLength negative",
			form: &Form{
				Fields: []Field{
					{
						Name: "text",
						Type: FieldTypeText,
						Validation: &Validation{
							MinLength: intPtr(-1),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot be negative",
		},
		{
			name: "min greater than max",
			form: &Form{
				Fields: []Field{
					{
						Name: "number",
						Type: FieldTypeNumber,
						Validation: &Validation{
							Min: floatPtr(100),
							Max: floatPtr(50),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot be greater than max",
		},
		{
			name: "minItems greater than maxItems",
			form: &Form{
				Fields: []Field{
					{
						Name: "items",
						Type: FieldTypeArray,
						Validation: &Validation{
							MinItems: intPtr(10),
							MaxItems: intPtr(5),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot be greater than maxItems",
		},
		{
			name: "step zero or negative",
			form: &Form{
				Fields: []Field{
					{
						Name: "number",
						Type: FieldTypeNumber,
						Validation: &Validation{
							Step: floatPtr(0),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "step must be positive",
		},
		{
			name: "text field with number validation",
			form: &Form{
				Fields: []Field{
					{
						Name: "text",
						Type: FieldTypeText,
						Validation: &Validation{
							Min: floatPtr(10),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "not applicable for field type",
		},
		{
			name: "number field with string validation",
			form: &Form{
				Fields: []Field{
					{
						Name: "number",
						Type: FieldTypeNumber,
						Validation: &Validation{
							MinLength: intPtr(5),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "not applicable for field type",
		},
		{
			name: "array field with string validation",
			form: &Form{
				Fields: []Field{
					{
						Name: "items",
						Type: FieldTypeArray,
						Validation: &Validation{
							MinLength: intPtr(5),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "not applicable for field type",
		},
		{
			name: "valid string validation",
			form: &Form{
				Fields: []Field{
					{
						Name: "text",
						Type: FieldTypeText,
						Validation: &Validation{
							MinLength: intPtr(5),
							MaxLength: intPtr(10),
							Pattern:   "^[A-Z]+$",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid number validation",
			form: &Form{
				Fields: []Field{
					{
						Name: "number",
						Type: FieldTypeNumber,
						Validation: &Validation{
							Min:  floatPtr(0),
							Max:  floatPtr(100),
							Step: floatPtr(0.5),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid array validation",
			form: &Form{
				Fields: []Field{
					{
						Name: "items",
						Type: FieldTypeArray,
						Validation: &Validation{
							MinItems: intPtr(1),
							MaxItems: intPtr(10),
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Form.Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestForm_Validate_ConditionalFields(t *testing.T) {
	tests := []struct {
		name    string
		form    *Form
		wantErr bool
		errMsg  string
	}{
		{
			name: "conditional field references non-existent field",
			form: &Form{
				Fields: []Field{
					{
						Name: "field1",
						Type: FieldTypeText,
						Conditional: &ConditionalField{
							Condition: "nonexistent",
							Then: []Field{
								{Name: "then_field", Type: FieldTypeText},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "references non-existent field",
		},
		{
			name: "conditional field with empty condition",
			form: &Form{
				Fields: []Field{
					{
						Name: "field1",
						Type: FieldTypeText,
						Conditional: &ConditionalField{
							Condition: "",
							Then: []Field{
								{Name: "then_field", Type: FieldTypeText},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "must specify a condition field name",
		},
		{
			name: "valid conditional field",
			form: &Form{
				Fields: []Field{
					{
						Name: "field1",
						Type: FieldTypeText,
					},
					{
						Name: "field2",
						Type: FieldTypeText,
						Conditional: &ConditionalField{
							Condition: "field1",
							Then: []Field{
								{Name: "then_field", Type: FieldTypeText},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "conditional field with duplicate names in then",
			form: &Form{
				Fields: []Field{
					{
						Name: "field1",
						Type: FieldTypeText,
					},
					{
						Name: "field2",
						Type: FieldTypeText,
						Conditional: &ConditionalField{
							Condition: "field1",
							Then: []Field{
								{Name: "duplicate", Type: FieldTypeText},
								{Name: "duplicate", Type: FieldTypeText},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate field name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Form.Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestForm_Validate_NestedFields(t *testing.T) {
	tests := []struct {
		name    string
		form    *Form
		wantErr bool
		errMsg  string
	}{
		{
			name: "duplicate field names in nested object",
			form: &Form{
				Fields: []Field{
					{
						Name: "user",
						Type: FieldTypeObject,
						Fields: []Field{
							{Name: "name", Type: FieldTypeText},
							{Name: "name", Type: FieldTypeEmail},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate field name",
		},
		{
			name: "valid nested object with unique fields",
			form: &Form{
				Fields: []Field{
					{
						Name: "user",
						Type: FieldTypeObject,
						Fields: []Field{
							{Name: "name", Type: FieldTypeText},
							{Name: "email", Type: FieldTypeEmail},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "deeply nested object",
			form: &Form{
				Fields: []Field{
					{
						Name: "user",
						Type: FieldTypeObject,
						Fields: []Field{
							{
								Name: "address",
								Type: FieldTypeObject,
								Fields: []Field{
									{Name: "street", Type: FieldTypeText},
									{Name: "city", Type: FieldTypeText},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "array field with nested item fields",
			form: &Form{
				Fields: []Field{
					{
						Name: "items",
						Type: FieldTypeArray,
						Fields: []Field{
							{Name: "item", Type: FieldTypeText},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Form.Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestForm_Validate_ComplexForm(t *testing.T) {
	// Test a complex valid form with multiple features
	form := &Form{
		Title:       "User Registration",
		Description: "Register a new user",
		Method:      "POST",
		Fields: []Field{
			{
				Name: "username",
				Type: FieldTypeText,
				Validation: &Validation{
					Required:  true,
					MinLength: intPtr(3),
					MaxLength: intPtr(20),
					Pattern:   "^[a-zA-Z0-9_]+$",
				},
			},
			{
				Name: "email",
				Type: FieldTypeEmail,
				Validation: &Validation{
					Required: true,
				},
			},
			{
				Name: "age",
				Type: FieldTypeNumber,
				Validation: &Validation{
					Min:  floatPtr(18),
					Max:  floatPtr(120),
					Step: floatPtr(1),
				},
			},
			{
				Name: "country",
				Type: FieldTypeSelect,
				Options: []Option{
					{Label: "USA", Value: "us"},
					{Label: "Canada", Value: "ca"},
					{Label: "UK", Value: "uk"},
				},
			},
			{
				Name: "address",
				Type: FieldTypeObject,
				Fields: []Field{
					{
						Name: "street",
						Type: FieldTypeText,
						Validation: &Validation{
							Required: true,
						},
					},
					{
						Name: "city",
						Type: FieldTypeText,
						Validation: &Validation{
							Required: true,
						},
					},
				},
			},
			{
				Name: "tags",
				Type: FieldTypeArray,
				Validation: &Validation{
					MinItems: intPtr(1),
					MaxItems: intPtr(10),
				},
				Fields: []Field{
					{Name: "tag", Type: FieldTypeText},
				},
			},
		},
	}

	err := form.Validate()
	if err != nil {
		t.Errorf("Form.Validate() error = %v, want nil for valid complex form", err)
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
