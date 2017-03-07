package native

import (
	"fmt"

	"github.com/pkg/errors"
)

type errDataGet struct {
	c           error
	entityType  string
	entityIDfmt string
	entityIDctx []interface{}
}

func dataGetError(cause error, typ string, idfmt string, ctx ...interface{}) error {
	if cause == nil {
		return nil
	}

	return errors.WithStack(&errDataGet{
		c:           cause,
		entityType:  typ,
		entityIDfmt: idfmt,
		entityIDctx: ctx,
	})
}

func (e *errDataGet) Error() string {
	id := fmt.Sprintf(e.entityIDfmt, e.entityIDctx...)
	return fmt.Sprintf("Error getting data for %v '%v': %v", e.entityType, id, e.c.Error())
}

func (e *errDataGet) Cause() error {
	return e.c
}
