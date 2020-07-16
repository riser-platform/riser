package e2e

import (
	"math/rand"
	"strings"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func ParseTestDummyEnv(envBody []byte) map[string]string {
	envMap := map[string]string{}
	lines := strings.Split(string(envBody), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}
