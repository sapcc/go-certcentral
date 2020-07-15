package go_certcentral

import "errors"

type Options struct {
	Token   string
	IsDebug bool
}

func (o *Options) validate() error {
	if o.Token == "" {
		return errors.New("token not provided")
	}

	return nil
}
