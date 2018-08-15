package runtime

import "github.com/brettbuddin/shaden/errors"

func typeRemainingError(typ string, arg int) error {
	return errors.Errorf("expects %s for remaining arguments (%d+)", typ, arg)
}
