package rand

import (
	"math"
	"os"
	"strconv"
	"time"
)

var seed uint32

func Seed() uint32 {
	seed = uint32(time.Now().UnixNano() + int64(os.Getpid()))

	return seed
}

func String(length int) string {
	randomNum := Seed()

	// 123 to 100...0123 then 000...0123 to standardize str length
	getStrForm := func(length int) string {
		blanketValue := uint32(math.Pow10(length + 1))

		return strconv.Itoa(int(blanketValue + randomNum%blanketValue))[1:]
	}

	return getStrForm(length)
}
