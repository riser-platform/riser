package e2e

import (
	"testing"
	"time"
)

// Step indicates a logical test Step. The message is printed in real time to stdout. Timing is printed upon completion
// I had too much friction w/Ginkgo and generally don't like strict BDD. This is trivial and good enough for real time output and timings.
// The big thing missing here is lack of structured output which can be fixed easily if we really need it. For now this is just for
// human consumption.
func Step(t *testing.T, message string, fn func()) {
	t.Helper()
	start := time.Now()
	fn()
	t.Logf("%s (%dms)\n", message, time.Since(start).Milliseconds())
}
