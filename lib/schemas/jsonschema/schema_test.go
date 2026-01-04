package jsonschema

import (
	"encoding/json"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*Schema) bool
	}{
		{
			name:    "valid simple schema",
			input:   `{"type": "string", "title": "Name"}`,
			wantErr: false,
			check: func(s *Schema) bool {
				return s.Title == "Name"
			},
		},
		{
			name:    "valid object schema with properties",
			input:   `{"type": "object", "properties": {"name": {"type": "string"}}}`,
			wantErr: false,
			check: func(s *Schema) bool {
				return s.Properties != nil && s.Properties["name"] != nil
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
			check:   nil,
		},
		{
			name:    "empty schema",
			input:   `{}`,
			wantErr: false,
			check: func(s *Schema) bool {
				return s != nil
			},
		},
		{
			name:    "schema with all core fields",
			input:   `{"$schema": "https://json-schema.org/draft/2020-12/schema", "$id": "test", "$ref": "#/definitions/test"}`,
			wantErr: false,
			check: func(s *Schema) bool {
				return s.Schema != "" && s.ID != "" && s.Ref != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(got) {
				t.Errorf("Parse() check failed")
			}
		})
	}
}

func TestSchema_GetType(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		wantStr  string
		wantArr  []string
		wantBool bool
	}{
		{
			name: "string type",
			schema: &Schema{
				Type: json.RawMessage(`"string"`),
			},
			wantStr:  "string",
			wantArr:  nil,
			wantBool: true,
		},
		{
			name: "array type",
			schema: &Schema{
				Type: json.RawMessage(`["string", "number"]`),
			},
			wantStr:  "",
			wantArr:  []string{"string", "number"},
			wantBool: true,
		},
		{
			name: "no type",
			schema: &Schema{
				Type: json.RawMessage(``),
			},
			wantStr:  "",
			wantArr:  nil,
			wantBool: false,
		},
		{
			name: "null type",
			schema: &Schema{
				Type: json.RawMessage(`"null"`),
			},
			wantStr:  "null",
			wantArr:  nil,
			wantBool: true,
		},
		{
			name: "number type",
			schema: &Schema{
				Type: json.RawMessage(`"number"`),
			},
			wantStr:  "number",
			wantArr:  nil,
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, gotArr, gotBool := tt.schema.GetType()
			if gotStr != tt.wantStr {
				t.Errorf("Schema.GetType() gotStr = %v, want %v", gotStr, tt.wantStr)
			}
			if len(gotArr) != len(tt.wantArr) {
				t.Errorf("Schema.GetType() gotArr length = %v, want %v", len(gotArr), len(tt.wantArr))
			} else {
				for i := range gotArr {
					if gotArr[i] != tt.wantArr[i] {
						t.Errorf("Schema.GetType() gotArr[%d] = %v, want %v", i, gotArr[i], tt.wantArr[i])
					}
				}
			}
			if gotBool != tt.wantBool {
				t.Errorf("Schema.GetType() gotBool = %v, want %v", gotBool, tt.wantBool)
			}
		})
	}
}
