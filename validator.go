package validator

import "context"

type Validator interface {
	Validate(ctx context.Context) error
}

type GroupValidator interface {
	Validate(ctx context.Context, validators []Validator) error
}

type simpleGroupValidator struct{}

func (v *simpleGroupValidator) Validate(ctx context.Context, validators []Validator) error {
	for _, v := range validators {
		if err := v.Validate(ctx); err != nil {
			return err
		}
	}
	return nil
}

func NewSimpleGroupValidator() *simpleGroupValidator {
	return &simpleGroupValidator{}
}

var _ (GroupValidator) = (*simpleGroupValidator)(nil)
