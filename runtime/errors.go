package runtime

import "buddin.us/shaden/errors"

func typeError(name, typ string, arg int) error {
	return errors.Errorf("%s expects %s for argument %d", name, typ, arg)
}

func typeRemainingError(name, typ string, arg int) error {
	return errors.Errorf("%s expects %s for remaining arguments (%d+)", name, typ, arg)
}

func exactArgCountError(name string, count int) error {
	return errors.Errorf("%s expects %d arguments", name, count)
}

func minArgCountError(name string, min int) error {
	return errors.Errorf("%s expects %d or more arguments", name, min)
}
