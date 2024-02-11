package controllers

import "errors"

var (
	ErrUnableToHitUrl  = errors.New("unable to hit the url")
	ErrInResponse      = errors.New("error in response")
	ErrInUnmarshalData = errors.New("error in unmarshal data")
)
