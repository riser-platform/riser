package e2e

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Step indicates a logical test Step. The message is printed in real time to stdout. Timing is printed upon completion
// I had too much friction w/Ginkgo and generally don't like strict BDD. This is trivial and good enough for real time output and timings.
// The big thing missing here is lack of structured output which can be fixed easily if we really need it. For now this is just for
// human consumption.
func Step(message string, fn func()) {
	fmt.Printf("• %s", message)
	start := time.Now()
	fn()
	fmt.Printf(" (%dms)\n", time.Since(start).Milliseconds())
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
