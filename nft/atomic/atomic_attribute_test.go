package atomic_test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/nft/atomic"
	"gotest.tools/assert"
)

func TestToAtomicAttribute(t *testing.T) {

	int8Val := int8(8)
	actual := atomic.ToAtomicAttribute(int8Val)
	expected := atomic.NewAtomicAttribute("int8", int8Val)
	assert.Assert(t, actual.IsEqual(expected))

	uint16Val := uint16(1)
	actual = atomic.ToAtomicAttribute(uint16Val)
	expected = atomic.NewAtomicAttribute("uint16", uint16Val)

	assert.Assert(t, actual.IsEqual(expected))

	uint64Val := uint64(10)
	actual = atomic.ToAtomicAttribute(uint64Val)
	expected = atomic.NewAtomicAttribute("uint64", uint64Val)
	assert.Assert(t, actual.IsEqual(expected))

	stringVal := "hola"
	actual = atomic.ToAtomicAttribute(stringVal)
	expected = atomic.NewAtomicAttribute("string", stringVal)

	assert.Assert(t, actual.IsEqual(expected))

	float32Val := float32(10.677)
	actual = atomic.ToAtomicAttribute(float32Val)
	expected = atomic.NewAtomicAttribute("float", float32Val)

	assert.Assert(t, actual.IsEqual(expected))

	float64Val := float64(10.677)
	actual = atomic.ToAtomicAttribute(float64Val)
	expected = atomic.NewAtomicAttribute("double", float64Val)

	assert.Assert(t, actual.IsEqual(expected))

	uint32Slice := []uint32{10, 11}
	actual = atomic.ToAtomicAttribute(uint32Slice)
	expected = atomic.NewAtomicAttribute("UINT32_VEC", uint32Slice)

	assert.Assert(t, actual.IsEqual(expected))

	float32Slice := []float32{10.12, 11.14}
	actual = atomic.ToAtomicAttribute(float32Slice)
	expected = atomic.NewAtomicAttribute("FLOAT_VEC", float32Slice)

	assert.Assert(t, actual.IsEqual(expected))

	float64Slice := []float64{10.12, 11.14}
	actual = atomic.ToAtomicAttribute(float64Slice)
	expected = atomic.NewAtomicAttribute("DOUBLE_VEC", float64Slice)

	assert.Assert(t, actual.IsEqual(expected))

}
