package utils

import (
	"strconv"
)

type (
	T interface{}
)

func Ternary[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}

	return falseVal
}

func StrToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
