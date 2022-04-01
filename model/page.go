package model

import "github.com/krogertechnology/krogo/pkg/errors"

type Page struct {
	Offset int `json:"offset"`
	Size   int `json:"size"`
}

// Check validation for pagination
func (p Page) Check() error {
	if p.Offset < 0 {
		return errors.InvalidParam{Param: []string{"page.offset"}}
	}

	if p.Size < 0 {
		return errors.InvalidParam{Param: []string{"page.size"}}
	}

	return nil
}
