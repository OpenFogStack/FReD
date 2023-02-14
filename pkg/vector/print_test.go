package vector

import (
	"testing"

	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/stretchr/testify/assert"
)

func TestSortedVCString(t *testing.T) {
	v := vclock.New()
	v["A"] = 1
	v["B"] = 40
	v["C"] = 0
	v["D"] = 10

	assert.Equal(t, "{\"A\":1, \"B\":40, \"C\":0, \"D\":10}", SortedVCString(v))
}
