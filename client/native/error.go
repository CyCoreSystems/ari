package native

import (
	"fmt"

	"github.com/rotisserie/eris"
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

	return eris.Wrap(&errDataGet{
		c:           cause,
		entityType:  typ,
		entityIDfmt: idfmt,
		entityIDctx: ctx,
	}, "failed to get data")
}

func (e *errDataGet) Error() string {
	id := fmt.Sprintf(e.entityIDfmt, e.entityIDctx...)
	return fmt.Sprintf("Error getting data for %v '%v': %v", e.entityType, id, e.c.Error())
}

func (e *errDataGet) Cause() error {
	return e.c
}
