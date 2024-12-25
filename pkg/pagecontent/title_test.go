package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTitleAuthorDate(t *testing.T) {
	testF := func(f string, title string, author string) {
		h, err := os.ReadFile(filepath.Join("testdata", f))
		assert.NoError(t, err)
		tad, err := ExtractTitleAuthorDate(string(h))
		assert.NoError(t, err)
		assert.Equal(t, title, tad.Title)
		assert.Equal(t, author, tad.Author)
	}

	testF("a.html", "Bluetooth Low Energy GATT Fuzzing", "Baptiste Boyer")
	testF("b.html", "HookCase：一款针对maxOS的逆向工程安全分析工具", "Alpha_h4ck")
	testF("c.html", "The Over-Engineering Pendulum", "Miłosz Smółka")
	testF("d.html", "Structured outputs", "")
	testF("e.html", `Optimising and Visualising Go Tests Parallelism: Why more cores don't speed up your Go tests`, "Robert Laszczak")

}
