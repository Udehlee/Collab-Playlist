package utils

import (
	"fmt"
	"math/rand"
)

func GenerateRandomState() string {
	return fmt.Sprintf("%d", rand.Intn(100000))
}
