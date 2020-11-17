package nasecache

import (
	"bytes"
	"encoding/gob"
	"unsafe"
)

// boolToByteArray serializes a bool as a byte array
func boolToByteArray(b bool) []byte {
	res := make([]byte, 1)
	if b {
		res[0] = 1
		return res
	}
	res[0] = 0
	return res
}

// byteArrayToBool deserializes a byte array to a bool
func byteArrayToBool(b []byte) bool {
	return b[0] == 1
}

// intToByteArray converts an int to a byte array
// uses unsafe operations
// inspired by https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func intToByteArray(num int) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

// byteArrayToInt converts a byte array to an int
// uses unsafe operations
// inspired by https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func byteArrayToInt(arr []byte) int {
	val := int(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

// boolToString converts bool to string
func boolToString(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

// genericToByteArray serializes an arbitrary struct into a byte array
// Caution: in should not be a pointer, while out must be a pointer to an empty variable of type []byte
func genericToByteArray(in interface{}, out *[]byte) error {
	buf := bytes.NewBuffer(*out)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(in)
	if err != nil {
		return err
	}
	*out = buf.Bytes()
	return nil
}

// byteArrayToGeneric deserializes byte arrays that were serialized by genericToByteArray
// Caution: out needs to be a pointer type of the type that was originally serialized or there will be errors
func byteArrayToGeneric(in []byte, out interface{}) error {
	buf := bytes.NewBuffer(in)
	dec := gob.NewDecoder(buf)

	return dec.Decode(out)
}
