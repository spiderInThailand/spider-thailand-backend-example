package random

import (
	"math/rand"
	"time"
)

type Random struct{}

var CHARTSET string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTYUVWXYZ123456789"

func NewRandom() *Random {
	return &Random{}
}

func (r *Random) RandomString(size int) string {

	ranByte := make([]byte, size)

	rand.Seed(time.Now().Unix())

	for i := 0; i < size; i++ {
		ranByte[i] = CHARTSET[rand.Intn(len(CHARTSET))]
	}

	return string(ranByte)

}
