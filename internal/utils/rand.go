package utils

import "github.com/pion/randutil"

const (
	RunesAlpha                 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	RunesDigit                 = "0123456789"
	runesCandidateIDFoundation = RunesAlpha + RunesDigit + "+/"

	lenUFrag = 16
	lenPwd   = 32
)

func GeneratePwd() (string, error) {
	return randutil.GenerateCryptoRandomString(lenPwd, RunesAlpha)
}

func GenerateUFrag() (string, error) {
	return randutil.GenerateCryptoRandomString(lenUFrag, RunesAlpha)
}
