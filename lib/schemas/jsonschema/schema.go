package jsonschema

import (
	"encoding/json"
)

// Parse unmarshals a JSON Schema string into a Schema struct
func Parse(schemaStr []byte) (*Schema, error) {
	var schema Schema
	err := json.Unmarshal(schemaStr, &schema)
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// Schema represents a JSON Schema document according to JSON Schema Draft 2020-12
type Schema struct {
	// Core vocabulary
	Schema        string             `json:"$schema,omitempty"`
	ID            string             `json:"$id,omitempty"`
	Anchor        string             `json:"$anchor,omitempty"`
	Ref           string             `json:"$ref,omitempty"`
	DynamicRef    string             `json:"$dynamicRef,omitempty"`
	DynamicAnchor string             `json:"$dynamicAnchor,omitempty"`
	Vocabularies  map[string]bool    `json:"$vocabularies,omitempty"`
	Comment       string             `json:"$comment,omitempty"`
	Defs          map[string]*Schema `json:"$defs,omitempty"`

	// Applicator vocabulary
	AllOf                []*Schema          `json:"allOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	Not                  *Schema            `json:"not,omitempty"`
	If                   *Schema            `json:"if,omitempty"`
	Then                 *Schema            `json:"then,omitempty"`
	Else                 *Schema            `json:"else,omitempty"`
	DependentSchemas     map[string]*Schema `json:"dependentSchemas,omitempty"`
	PrefixItems          []*Schema          `json:"prefixItems,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	Contains             *Schema            `json:"contains,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	PatternProperties    map[string]*Schema `json:"patternProperties,omitempty"`
	AdditionalProperties *Schema            `json:"additionalProperties,omitempty"`
	PropertyNames        *Schema            `json:"propertyNames,omitempty"`

	// Unevaluated vocabulary
	UnevaluatedItems      *Schema `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties *Schema `json:"unevaluatedProperties,omitempty"`

	// Validation vocabulary - type and content
	Type  json.RawMessage `json:"type,omitempty"` // Can be string or []string
	Enum  []any           `json:"enum,omitempty"`
	Const any             `json:"const,omitempty"`

	// Validation vocabulary - numbers
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`

	// Validation vocabulary - strings
	MaxLength        *int    `json:"maxLength,omitempty"`
	MinLength        *int    `json:"minLength,omitempty"`
	Pattern          string  `json:"pattern,omitempty"`
	Format           string  `json:"format,omitempty"`
	ContentEncoding  string  `json:"contentEncoding,omitempty"`
	ContentMediaType string  `json:"contentMediaType,omitempty"`
	ContentSchema    *Schema `json:"contentSchema,omitempty"`

	// Validation vocabulary - arrays
	MaxItems    *int  `json:"maxItems,omitempty"`
	MinItems    *int  `json:"minItems,omitempty"`
	UniqueItems *bool `json:"uniqueItems,omitempty"`
	MaxContains *int  `json:"maxContains,omitempty"`
	MinContains *int  `json:"minContains,omitempty"`

	// Validation vocabulary - objects
	MaxProperties     *int                `json:"maxProperties,omitempty"`
	MinProperties     *int                `json:"minProperties,omitempty"`
	Required          []string            `json:"required,omitempty"`
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"`

	// Meta-data vocabulary
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Deprecated  *bool  `json:"deprecated,omitempty"`
	ReadOnly    *bool  `json:"readOnly,omitempty"`
}

// GetType returns the type as a string or slice of strings
// This helper method handles the fact that "type" can be either a string or an array of strings
func (s *Schema) GetType() (string, []string, bool) {
	if len(s.Type) == 0 {
		return "", nil, false
	}

	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(s.Type, &str); err == nil {
		return str, nil, true
	}

	// Try to unmarshal as array of strings
	var arr []string
	if err := json.Unmarshal(s.Type, &arr); err == nil {
		return "", arr, true
	}

	return "", nil, false
}
