package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTitleAuthorDate(t *testing.T) {
	testF := func(f string, s string) {
		h, err := os.ReadFile(filepath.Join("testdata", f))
		assert.NoError(t, err)
		tad, err := ExtractTitleAuthorDate(string(h))
		assert.NoError(t, err)
		assert.Equal(t, s, tad.Title)
	}

	testF("a.html", "Bluetooth Low Energy GATT Fuzzing")
	testF("b.html", "HookCase：一款针对maxOS的逆向工程安全分析工具")
	testF("c.html", "The Over-Engineering Pendulum")
	testF("d.html", "Structured outputs")
	testF("e.html", `Optimising and Visualising Go Tests Parallelism: Why more cores don't speed up your Go tests`)

}
