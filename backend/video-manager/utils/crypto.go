package utils

import (
	"crypto/rand"
	"math/big"
)

// Genrates a random string of length 10 for videoId
func GenerateVideoId() string {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length := 10
	max := big.NewInt(62)
	rstr := ""

	for i := 1; i <= length; i++ {
		secureNum, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}

		rstr += string(alphabet[int(secureNum.Int64())])
	}
	return rstr
}
