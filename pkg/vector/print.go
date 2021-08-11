package vector

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/DistributedClocks/GoVector/govec/vclock"
)

// SortedVCString returns a deterministic string encoding of a vector clock
func SortedVCString(v vclock.VClock) string {
	//sort

	ids := make([]string, len(v))
	i := 0
	for id := range v {
		ids[i] = id
		i++
	}

	sort.Strings(ids)

	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i := range ids {
		buffer.WriteString(fmt.Sprintf("\"%s\":%d", ids[i], v[ids[i]]))
		if i+1 < len(ids) {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}
