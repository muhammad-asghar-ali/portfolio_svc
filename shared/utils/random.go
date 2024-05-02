package utils

import (
	"math/rand"
	"time"

	"github.com/0xbase-Corp/portfolio_svc/shared/configs"
)

var r *rand.Rand

func init() {
	s := rand.NewSource(time.Now().UTC().UnixNano())
	r = rand.New(s)
}

func RandomString(n int) string {
	letterRunes := []rune(configs.EnvConfigVars.GetSecret())

	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}

	return string(b)
}
