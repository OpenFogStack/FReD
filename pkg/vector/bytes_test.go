package vector

import (
	"testing"

	"git.tu-berlin.de/mcc-fred/vclock"
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

func TestFromBytes(t *testing.T) {
	v := vclock.New()
	v["A"] = 1
	v["B"] = 40
	v["C"] = 0
	v["D"] = 10

	b := Bytes(v)

	b[1] = b[1] ^ 0x11

	c, err := FromBytes(b)

	assert.Error(t, err)
	assert.False(t, v.Compare(c, vclock.Equal))
}
