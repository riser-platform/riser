package e2e

import (
	"sync"
	"testing"
	"time"

	"github.com/wzshiming/ctc"
)

// Step indicates a logical test Step. The message is printed in real time to stdout. Timing is printed upon completion
// I had too much friction w/Ginkgo and generally don't like strict BDD. This is trivial and good enough for real time output and timings.
// The big thing missing here is lack of structured output which can be fixed easily if we really need it. For now this is just for
// human consumption.
func Step(t *testing.T, message string, fn func()) {
	t.Helper()
	start := time.Now()
	fn()
	t.Logf("%s%s (%dms)%s\n", stepColor(t.Name()), message, time.Since(start).Milliseconds(), ctc.Reset)
}

var stepColors = map[string]ctc.Color{}
var allStepColors = []ctc.Color{
	ctc.ForegroundBrightBlue,
	ctc.ForegroundCyan,
	ctc.ForegroundMagenta,
	ctc.ForegroundYellow,
	ctc.ForegroundGreen,
	ctc.ForegroundBrightBlack,
}
var stepColorIdx int
var stepColorMux sync.Mutex

func stepColor(testName string) ctc.Color {
	stepColorMux.Lock()
	defer stepColorMux.Unlock()

	color, ok := stepColors[testName]
	if ok {
		return color
	}

	color = allStepColors[stepColorIdx]
	stepColorIdx = (stepColorIdx + 1) % len(allStepColors)
	stepColors[testName] = color
	return color
}
