package nft_test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/nft"
	"gotest.tools/assert"
)

func TestToAtomicAttribute(t *testing.T) {

	int8Val := int8(8)
	actual := nft.ToAtomicAttribute(int8Val)
	expected := nft.NewAtomicAttribute("int8", int8Val)
	assert.Assert(t, actual.IsEqual(expected))

	uint16Val := uint16(1)
	actual = nft.ToAtomicAttribute(uint16Val)
	expected = nft.NewAtomicAttribute("uint16", uint16Val)

	assert.Assert(t, actual.IsEqual(expected))

	uint64Val := uint64(10)
	actual = nft.ToAtomicAttribute(uint64Val)
	expected = nft.NewAtomicAttribute("uint64", uint64Val)
	assert.Assert(t, actual.IsEqual(expected))

	stringVal := "hola"
	actual = nft.ToAtomicAttribute(stringVal)
	expected = nft.NewAtomicAttribute("string", stringVal)

	assert.Assert(t, actual.IsEqual(expected))

	float32Val := float32(10.677)
	actual = nft.ToAtomicAttribute(float32Val)
	expected = nft.NewAtomicAttribute("float", float32Val)

	assert.Assert(t, actual.IsEqual(expected))

	float64Val := float64(10.677)
	actual = nft.ToAtomicAttribute(float64Val)
	expected = nft.NewAtomicAttribute("double", float64Val)

	assert.Assert(t, actual.IsEqual(expected))

	uint32Slice := []uint32{10, 11}
	actual = nft.ToAtomicAttribute(uint32Slice)
	expected = nft.NewAtomicAttribute("UINT32_VEC", uint32Slice)

	assert.Assert(t, actual.IsEqual(expected))

	float32Slice := []float32{10.12, 11.14}
	actual = nft.ToAtomicAttribute(float32Slice)
	expected = nft.NewAtomicAttribute("FLOAT_VEC", float32Slice)

	assert.Assert(t, actual.IsEqual(expected))

	float64Slice := []float64{10.12, 11.14}
	actual = nft.ToAtomicAttribute(float64Slice)
	expected = nft.NewAtomicAttribute("DOUBLE_VEC", float64Slice)

	assert.Assert(t, actual.IsEqual(expected))

}
