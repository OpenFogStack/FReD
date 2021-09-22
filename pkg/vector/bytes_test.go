package vector

import (
	"testing"

	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	v := vclock.New()
	v["A"] = 1
	v["B"] = 40
	v["C"] = 0
	v["D"] = 10

	c, err := FromBytes(Bytes(v))

	assert.NoError(t, err)
	assert.True(t, v.Compare(c, vclock.Equal))
}