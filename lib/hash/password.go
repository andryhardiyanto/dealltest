package hash

import (
	"encoding/base64"

	"github.com/AndryHardiyanto/dealltest/lib/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func ComparePassword(pwd, value string) bool {
	pwdBytes, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}

	err = bcrypt.CompareHashAndPassword(pwdBytes, []byte(value))
	if err != nil {
		log.Error().Msg(err.Error())
		return false
	}

	return true
}
