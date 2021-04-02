package generator

import (
	"kudago/pkg/constants"
	"math/rand"
)

func RandStringRunes(n uint8) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = constants.LetterRunes[rand.Intn(len(constants.LetterRunes))]
	}
	return string(b)
}
