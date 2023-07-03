package helpers

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomOTP(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	min := int64(Pow(10, length-1))
	max := int64(Pow(10, length) - 1)

	otp := seededRand.Int63n(max-min+1) + min

	return strconv.Itoa(int(otp))
}

func Pow(x, y int) int {
	if y == 0 {
		return 1
	}

	result := x
	for i := 1; i < y; i++ {
		result *= x
	}

	return result
}
