package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTitleAuthorDate(t *testing.T) {
	testF := func(f string, title string, author string, date string) {
		h, err := os.ReadFile(filepath.Join("testdata", f))
		assert.NoError(t, err)
		tad, err := ExtractTitleAuthorDate(string(h))
		assert.NoError(t, err)
		assert.Equal(t, title, tad.Title)
		assert.Equal(t, author, tad.Author)
		assert.Equal(t, date, tad.Date)
	}

	testF("a.html", "Bluetooth Low Energy GATT Fuzzing", "Baptiste Boyer", "2024-10-25")
	testF("b.html", "HookCase：一款针对maxOS的逆向工程安全分析工具", "Alpha_h4ck", "2024-12-09")
	testF("c.html", "The Over-Engineering Pendulum", "Miłosz Smółka", "2024-12-17")
	testF("d.html", "Structured outputs", "", "2024-12-06")
	testF("e.html", `Optimising and Visualising Go Tests Parallelism: Why more cores don't speed up your Go tests`, "Robert Laszczak", "2024-10-17")

}
