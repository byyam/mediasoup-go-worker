package utils

import "github.com/pion/randutil"

const (
	runesAlpha                 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runesDigit                 = "0123456789"
	runesCandidateIDFoundation = runesAlpha + runesDigit + "+/"

	lenUFrag = 16
	lenPwd   = 32
)

func GeneratePwd() (string, error) {
	return randutil.GenerateCryptoRandomString(lenPwd, runesAlpha)
}

func GenerateUFrag() (string, error) {
	return randutil.GenerateCryptoRandomString(lenUFrag, runesAlpha)
}
