package html

import (
	"context"
	"io"

	"github.com/Olian04/form-from-schema/lib"
)

func ConvertFormToHtml(ctx context.Context, form *lib.Form, w io.Writer) error {
	return Form(form).Render(ctx, w)
}
