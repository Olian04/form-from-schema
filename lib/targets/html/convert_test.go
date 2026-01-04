package html

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Olian04/form-from-schema/lib"
)

func TestConvertFormToHtml(t *testing.T) {
	tests := []struct {
		name         string
		form         *lib.Form
		wantErr      bool
		wantContains []string
		notContains  []string
	}{
		{
			name:    "nil form causes panic",
			form:    nil,
			wantErr: true, // templ doesn't handle nil gracefully - will panic
		},
		{
			name: "empty form",
			form: &lib.Form{
				Fields: []lib.Field{},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				"</form>",
			},
		},
		{
			name: "simple form with title and description",
			form: &lib.Form{
				Title:       "Test Form",
				Description: "This is a test form",
				Method:      "POST",
				Action:      "/submit",
				Fields: []lib.Field{
					{
						Name: "username",
						Type: lib.FieldTypeText,
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				`method="POST"`,
				`action="/submit"`,
				"<h1>Test Form</h1>",
				"This is a test form",
				"</form>",
			},
		},
		{
			name: "form with text field",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name:        "username",
						Type:        lib.FieldTypeText,
						Label:       "Username",
						Description: "Enter your username",
						Placeholder: "username",
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				"username",
				"</form>",
			},
		},
		{
			name: "form with multiple fields",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "username",
						Type: lib.FieldTypeText,
					},
					{
						Name: "email",
						Type: lib.FieldTypeEmail,
					},
					{
						Name: "age",
						Type: lib.FieldTypeNumber,
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				"username",
				"email",
				"age",
				"</form>",
			},
		},
		{
			name: "form with select field",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "country",
						Type: lib.FieldTypeSelect,
						Options: []lib.Option{
							{Label: "USA", Value: "us"},
							{Label: "Canada", Value: "ca"},
						},
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				"country",
				"</form>",
			},
		},
		{
			name: "form with nested object field",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "user",
						Type: lib.FieldTypeObject,
						Fields: []lib.Field{
							{
								Name: "name",
								Type: lib.FieldTypeText,
							},
							{
								Name: "email",
								Type: lib.FieldTypeEmail,
							},
						},
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				"</form>",
			},
		},
		{
			name: "form with array field",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "tags",
						Type: lib.FieldTypeArray,
						Fields: []lib.Field{
							{
								Name: "tag",
								Type: lib.FieldTypeText,
							},
						},
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				"</form>",
			},
		},
		{
			name: "form with submit button",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "field",
						Type: lib.FieldTypeText,
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				`<button type="submit"`,
				"Submit",
				"</button>",
			},
		},
		{
			name: "form with all attributes",
			form: &lib.Form{
				Title:       "Registration Form",
				Description: "Please fill out this form",
				Method:      "POST",
				Action:      "/register",
				Fields: []lib.Field{
					{
						Name:        "username",
						Type:        lib.FieldTypeText,
						Label:       "Username",
						Description: "Choose a username",
						Placeholder: "Enter username",
						Default:     "guest",
					},
					{
						Name: "email",
						Type: lib.FieldTypeEmail,
						Validation: &lib.Validation{
							Required: true,
						},
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				`method="POST"`,
				`action="/register"`,
				"<h1>Registration Form</h1>",
				"Please fill out this form",
				"</form>",
			},
		},
		{
			name: "form without method defaults to empty",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "field",
						Type: lib.FieldTypeText,
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				`method=""`,
			},
		},
		{
			name: "form without action defaults to empty",
			form: &lib.Form{
				Fields: []lib.Field{
					{
						Name: "field",
						Type: lib.FieldTypeText,
					},
				},
			},
			wantErr: false,
			wantContains: []string{
				"<form",
				`action=""`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var buf bytes.Buffer

			// Handle nil form case separately as it causes panic
			if tt.form == nil {
				defer func() {
					if r := recover(); r == nil && !tt.wantErr {
						t.Error("ConvertFormToHtml() with nil form should panic")
					}
				}()
			}

			err := ConvertFormToHtml(ctx, tt.form, &buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertFormToHtml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			output := buf.String()

			// Check that expected strings are contained
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("ConvertFormToHtml() output does not contain %q. Output: %s", want, output)
				}
			}

			// Check that unexpected strings are not contained
			for _, notWant := range tt.notContains {
				if strings.Contains(output, notWant) {
					t.Errorf("ConvertFormToHtml() output contains unexpected %q. Output: %s", notWant, output)
				}
			}
		})
	}
}

func TestConvertFormToHtml_NilWriter(t *testing.T) {
	form := &lib.Form{
		Fields: []lib.Field{
			{
				Name: "test",
				Type: lib.FieldTypeText,
			},
		},
	}

	ctx := context.Background()

	// nil writer causes panic in templ
	defer func() {
		if r := recover(); r == nil {
			t.Error("ConvertFormToHtml() with nil writer should panic")
		}
	}()

	_ = ConvertFormToHtml(ctx, form, nil)
}

func TestConvertFormToHtml_AllFieldTypes(t *testing.T) {
	fieldTypes := []lib.FieldType{
		lib.FieldTypeText,
		lib.FieldTypeEmail,
		lib.FieldTypePassword,
		lib.FieldTypeNumber,
		lib.FieldTypeTel,
		lib.FieldTypeURL,
		lib.FieldTypeDate,
		lib.FieldTypeTime,
		lib.FieldTypeDateTime,
		lib.FieldTypeMonth,
		lib.FieldTypeWeek,
		lib.FieldTypeTextarea,
		lib.FieldTypeSelect,
		lib.FieldTypeCheckbox,
		lib.FieldTypeRadio,
		lib.FieldTypeFile,
		lib.FieldTypeHidden,
		lib.FieldTypeObject,
		lib.FieldTypeArray,
	}

	for _, fieldType := range fieldTypes {
		t.Run(string(fieldType), func(t *testing.T) {
			form := &lib.Form{
				Fields: []lib.Field{
					{
						Name: "test_field",
						Type: fieldType,
					},
				},
			}

			ctx := context.Background()
			var buf bytes.Buffer

			err := ConvertFormToHtml(ctx, form, &buf)
			if err != nil {
				t.Errorf("ConvertFormToHtml() error = %v for field type %s", err, fieldType)
				return
			}

			output := buf.String()
			if !strings.Contains(output, "<form") {
				t.Errorf("ConvertFormToHtml() output does not contain form tag for field type %s", fieldType)
			}
		})
	}
}

func TestConvertFormToHtml_ComplexForm(t *testing.T) {
	form := &lib.Form{
		Title:       "User Registration",
		Description: "Register a new user account",
		Method:      "POST",
		Action:      "/register",
		Fields: []lib.Field{
			{
				Name:        "username",
				Type:        lib.FieldTypeText,
				Label:       "Username",
				Description: "Choose a unique username",
				Placeholder: "username",
				Default:     "",
				Validation: &lib.Validation{
					Required:  true,
					MinLength: intPtr(3),
					MaxLength: intPtr(20),
					Pattern:   "^[a-zA-Z0-9_]+$",
				},
			},
			{
				Name:  "email",
				Type:  lib.FieldTypeEmail,
				Label: "Email Address",
				Validation: &lib.Validation{
					Required: true,
				},
			},
			{
				Name:  "age",
				Type:  lib.FieldTypeNumber,
				Label: "Age",
				Validation: &lib.Validation{
					Min:  floatPtr(18),
					Max:  floatPtr(120),
					Step: floatPtr(1),
				},
			},
			{
				Name:  "country",
				Type:  lib.FieldTypeSelect,
				Label: "Country",
				Options: []lib.Option{
					{Label: "United States", Value: "us"},
					{Label: "Canada", Value: "ca"},
					{Label: "United Kingdom", Value: "uk"},
				},
			},
			{
				Name:  "address",
				Type:  lib.FieldTypeObject,
				Label: "Address",
				Fields: []lib.Field{
					{
						Name:  "street",
						Type:  lib.FieldTypeText,
						Label: "Street",
					},
					{
						Name:  "city",
						Type:  lib.FieldTypeText,
						Label: "City",
					},
					{
						Name:  "zip",
						Type:  lib.FieldTypeText,
						Label: "ZIP Code",
					},
				},
			},
			{
				Name:  "tags",
				Type:  lib.FieldTypeArray,
				Label: "Tags",
				Fields: []lib.Field{
					{
						Name: "tag",
						Type: lib.FieldTypeText,
					},
				},
				Validation: &lib.Validation{
					MinItems: intPtr(1),
					MaxItems: intPtr(10),
				},
			},
		},
	}

	ctx := context.Background()
	var buf bytes.Buffer

	err := ConvertFormToHtml(ctx, form, &buf)
	if err != nil {
		t.Fatalf("ConvertFormToHtml() error = %v", err)
	}

	output := buf.String()

	// Verify form structure
	expectedElements := []string{
		"<form",
		`method="POST"`,
		`action="/register"`,
		"<h1>User Registration</h1>",
		"Register a new user account",
		"</form>",
		`<button type="submit"`,
		"Submit",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(output, elem) {
			t.Errorf("ConvertFormToHtml() output missing expected element: %q", elem)
		}
	}

	// Verify output is valid HTML structure
	if !strings.HasPrefix(strings.TrimSpace(output), "<form") {
		t.Error("ConvertFormToHtml() output should start with <form tag")
	}

	if !strings.Contains(output, "</form>") {
		t.Error("ConvertFormToHtml() output should contain closing </form> tag")
	}
}

func TestConvertFormToHtml_ContextCancellation(t *testing.T) {
	form := &lib.Form{
		Fields: []lib.Field{
			{
				Name: "test",
				Type: lib.FieldTypeText,
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var buf bytes.Buffer
	err := ConvertFormToHtml(ctx, form, &buf)

	// Context cancellation might or might not cause an error depending on templ implementation
	// We just verify the function handles it gracefully
	if err != nil {
		t.Logf("ConvertFormToHtml() returned error on cancelled context (expected): %v", err)
	}
}

func TestConvertFormToHtml_EmptyFields(t *testing.T) {
	form := &lib.Form{
		Title:  "Empty Form",
		Fields: []lib.Field{},
	}

	ctx := context.Background()
	var buf bytes.Buffer

	err := ConvertFormToHtml(ctx, form, &buf)
	if err != nil {
		t.Errorf("ConvertFormToHtml() error = %v, want nil for form with empty fields", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<form") {
		t.Error("ConvertFormToHtml() should generate form tag even with empty fields")
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
