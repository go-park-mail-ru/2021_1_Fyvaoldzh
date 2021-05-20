package generator

import (
	"crypto/sha256"
	"encoding/base64"
	"kudago/pkg/constants"
	"math/rand"
	"net/http"
	"time"
)

func RandStringRunes(n uint8) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = constants.LetterRunes[rand.Intn(len(constants.LetterRunes))]
	}
	return string(b)
}

func HashPassword(oldPassword string) string {
	hash := sha256.New()
	salt := RandStringRunes(constants.SaltLength)
	hash.Write([]byte(salt + oldPassword))
	return salt + base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func CheckHashedPassword(databasePassword string, gotPassword string) bool {
	salt := databasePassword[:8]
	hash := sha256.New()
	hash.Write([]byte(salt + gotPassword))
	gotPassword = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	return gotPassword == databasePassword[8:]
}

func CreateCookieValue(n uint8) string {
	key := RandStringRunes(n)
	return key
}

func CreateCookieWithValue(value string) *http.Cookie {
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    value,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
	}

	return newCookie
}
func CreateCookie(n uint8) *http.Cookie {
	key := RandStringRunes(n)

	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}

	return newCookie
}
