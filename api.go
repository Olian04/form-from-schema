package formfromschema

import (
	"context"
	"io"

	"github.com/Olian04/form-from-schema/lib"
	"github.com/Olian04/form-from-schema/lib/schemas/jsonschema"
	"github.com/Olian04/form-from-schema/lib/targets/html"
)

// FromJsonSchema parses a JSON Schema and converts it to a Form struct
// The Form struct is NOT validated and should be validated before use by the caller
func FromJsonSchema(schema []byte) (*lib.Form, error) {
	jsonSchema, err := jsonschema.Parse(schema)
	if err != nil {
		return nil, err
	}
	form, err := jsonschema.ConvertSchemaToForm(jsonSchema)
	if err != nil {
		return nil, err
	}
	return form, nil
}

// ToHtml converts a Form struct to HTML and writes it to the provided writer
func ToHtml(ctx context.Context, form *lib.Form, w io.Writer) error {
	return html.ConvertFormToHtml(ctx, form, w)
}
