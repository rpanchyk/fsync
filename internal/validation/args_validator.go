package validation

import "errors"

// TODO: remove
type ArgsValidator struct {
}

func (v *ArgsValidator) Validate(args []string) error {
	if len(args) != 2 {
		return errors.New("invalid args")
	}
	return nil
}
