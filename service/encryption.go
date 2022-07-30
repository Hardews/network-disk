package service

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func Encryption(password string) (error, string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return err, string(hash)
	}
	return err, string(hash)

}

func Interpretation(passwordInSql, passwordInput string) (error, bool) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordInSql), []byte(passwordInput)) //验证（对比）
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = nil
			return err, false
		}
		return err, false
	} else {
		return err, true
	}
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWSYZabcdefghijklmnopqrstuvwsyz1234567890")

func RandomStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
