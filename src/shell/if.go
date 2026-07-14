package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

type If string

const booleanTrue = "true"

func (i If) Ignore() bool {
	if len(i) == 0 {
		return false
	}

	got, err := evaluate(i)
	if err != nil {
		return false
	}

	return got == booleanTrue
}

func (i If) Validate() error {
	if len(i) == 0 {
		return nil
	}

	_, err := evaluate(i)
	return err
}

func evaluate(i If) (string, error) {
	template := fmt.Sprintf(`{{ if %s }}false{{ else }}true{{ end }}`, i)
	ctx := currentRuntime()
	if ctx == nil {
		ctx = &context.Runtime{
			Shell: BASH,
			OS:    context.LINUX,
			Home:  context.Home(),
		}
	}

	return parse(template, ctx)
}
