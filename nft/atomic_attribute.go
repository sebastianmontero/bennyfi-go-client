package nft

import (
	"fmt"
	"log"
	"strings"

	"github.com/eoscanada/eos-go"
)

var AtomicAttributeVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "int8", Type: int8(0)},
	{Name: "int16", Type: int16(0)},
	{Name: "int32", Type: int32(0)},
	{Name: "int64", Type: int64(0)},
	{Name: "uint8", Type: uint8(0)},
	{Name: "uint16", Type: uint16(0)},
	{Name: "uint32", Type: uint32(0)},
	{Name: "uint64", Type: uint64(0)},
	{Name: "float", Type: float32(0)},
	{Name: "double", Type: float64(0)},
	{Name: "string", Type: ""},
	{Name: "INT8_VEC", Type: []int8{}},
	{Name: "INT16_VEC", Type: []int16{}},
	{Name: "INT32_VEC", Type: []int32{}},
	{Name: "INT64_VEC", Type: []int64{}},
	{Name: "UINT8_VEC", Type: []uint8{}},
	{Name: "UINT16_VEC", Type: []uint16{}},
	{Name: "UINT32_VEC", Type: []uint32{}},
	{Name: "UINT64_VEC", Type: []uint64{}},
	{Name: "FLOAT_VEC", Type: []float32{}},
	{Name: "DOUBLE_VEC", Type: []float64{}},
	{Name: "STRING_VEC", Type: []string{}},
})

type InvalidTypeError struct {
	Label        string
	ExpectedType string
	Attribute    *AtomicAttribute
}

func (c *InvalidTypeError) Error() string {
	return fmt.Sprintf("received an unexpected type %T for metadata variant %T", c.ExpectedType, c.Attribute)
}

type AttributeMap map[string]*AtomicAttribute

type AtomicAttribute struct {
	*eos.BaseVariant
}

func NewAtomicAttribute(typeId string, value interface{}) *AtomicAttribute {
	return &AtomicAttribute{
		&eos.BaseVariant{
			TypeID: AtomicAttributeVariant.TypeID(typeId),
			Impl:   value,
		},
	}
}

func ToAtomicAttribute(value interface{}) *AtomicAttribute {
	switch v := value.(type) {
	case float32:
		return NewAtomicAttribute("float", v)
	case float64:
		return NewAtomicAttribute("double", v)
	case []int8, []int16, []int32, []int64, []uint8, []uint16, []uint32, []uint64, []string:
		typeId := fmt.Sprintf("%v_VEC", strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%T", v), "[]", "")))
		return NewAtomicAttribute(typeId, v)
	case []float32:
		return NewAtomicAttribute("FLOAT_VEC", v)
	case []float64:
		return NewAtomicAttribute("DOUBLE_VEC", v)
	default:
		return NewAtomicAttribute(fmt.Sprintf("%T", v), v)
	}
}

func (m *AtomicAttribute) InvalidTypeError(expectedType string) *InvalidTypeError {
	return &InvalidTypeError{
		Label:        fmt.Sprintf("received an unexpected type %T for variant %T", m.Impl, m),
		ExpectedType: "int8",
		Attribute:    m,
	}
}

func (m *AtomicAttribute) String() string {
	return fmt.Sprint(m.Impl)
}

func (m *AtomicAttribute) Int8() (int8, error) {
	switch v := m.Impl.(type) {
	case int8:
		return v, nil
	default:
		return 0, m.InvalidTypeError("int8")
	}
}

func (m *AtomicAttribute) Int16() (int16, error) {
	switch v := m.Impl.(type) {
	case int16:
		return v, nil
	default:
		return 0, m.InvalidTypeError("int16")
	}
}

func (m *AtomicAttribute) Int32() (int32, error) {
	switch v := m.Impl.(type) {
	case int32:
		return v, nil
	default:
		return 0, m.InvalidTypeError("int32")
	}
}

func (m *AtomicAttribute) Int64() (int64, error) {
	switch v := m.Impl.(type) {
	case int64:
		return v, nil
	default:
		return 0, m.InvalidTypeError("int64")
	}
}

func (m *AtomicAttribute) UInt8() (uint8, error) {
	switch v := m.Impl.(type) {
	case uint8:
		return v, nil
	default:
		return 0, m.InvalidTypeError("uint8")
	}
}

func (m *AtomicAttribute) UInt16() (uint16, error) {
	switch v := m.Impl.(type) {
	case uint16:
		return v, nil
	default:
		return 0, m.InvalidTypeError("uint16")
	}
}

func (m *AtomicAttribute) UInt32() (uint32, error) {
	switch v := m.Impl.(type) {
	case uint32:
		return v, nil
	default:
		return 0, m.InvalidTypeError("uint32")
	}
}

func (m *AtomicAttribute) UInt64() (uint64, error) {
	switch v := m.Impl.(type) {
	case uint64:
		return v, nil
	default:
		return 0, m.InvalidTypeError("uint64")
	}
}

func (m *AtomicAttribute) Float32() (float32, error) {
	switch v := m.Impl.(type) {
	case float32:
		return v, nil
	default:
		return 0, m.InvalidTypeError("float32")
	}
}

func (m *AtomicAttribute) Float64() (float64, error) {
	switch v := m.Impl.(type) {
	case float64:
		return v, nil
	default:
		return 0, m.InvalidTypeError("float64")
	}
}

func (m *AtomicAttribute) Int8Slice() ([]int8, error) {
	switch v := m.Impl.(type) {
	case []int8:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]int8")
	}
}

func (m *AtomicAttribute) Int16Slice() ([]int16, error) {
	switch v := m.Impl.(type) {
	case []int16:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]int16")
	}
}

func (m *AtomicAttribute) Int32Slice() ([]int32, error) {
	switch v := m.Impl.(type) {
	case []int32:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]int32")
	}
}

func (m *AtomicAttribute) Int64Slice() ([]int64, error) {
	switch v := m.Impl.(type) {
	case []int64:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]int64")
	}
}

func (m *AtomicAttribute) UInt8Slice() ([]uint8, error) {
	switch v := m.Impl.(type) {
	case []uint8:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]uint8")
	}
}

func (m *AtomicAttribute) UInt16Slice() ([]uint16, error) {
	switch v := m.Impl.(type) {
	case []uint16:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]uint16")
	}
}

func (m *AtomicAttribute) UInt32Slice() ([]uint32, error) {
	switch v := m.Impl.(type) {
	case []uint32:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]uint32")
	}
}

func (m *AtomicAttribute) UInt64Slice() ([]uint64, error) {
	switch v := m.Impl.(type) {
	case []uint64:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]uint64")
	}
}

func (m *AtomicAttribute) Float32Slice() ([]float32, error) {
	switch v := m.Impl.(type) {
	case []float32:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]float32")
	}
}

func (m *AtomicAttribute) Float64Slice() ([]float64, error) {
	switch v := m.Impl.(type) {
	case []float64:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]float64")
	}
}

func (m *AtomicAttribute) StringSlice() ([]string, error) {
	switch v := m.Impl.(type) {
	case []string:
		return v, nil
	default:
		return nil, m.InvalidTypeError("[]string")
	}
}

// IsEqual evaluates if the two AtomicAttributes have the same types and values (deep compare)
func (m *AtomicAttribute) IsEqual(m2 *AtomicAttribute) bool {

	if m.TypeID != m2.TypeID {
		log.Println("AtomicAttribute types inequal: ", m.TypeID, " vs ", m2.TypeID)
		return false
	}

	if m.String() != m2.String() {
		log.Println("AtomicAttribute Values.String() inequal: ", m.String(), " vs ", m2.String())
		return false
	}

	return true
}

// MarshalJSON translates to []byte
func (m *AtomicAttribute) MarshalJSON() ([]byte, error) {
	return m.BaseVariant.MarshalJSON(AtomicAttributeVariant)
}

// UnmarshalJSON translates AtomicAttributeVariant
func (m *AtomicAttribute) UnmarshalJSON(data []byte) error {
	return m.BaseVariant.UnmarshalJSON(data, AtomicAttributeVariant)
}

// UnmarshalBinary ...
func (m *AtomicAttribute) UnmarshalBinary(decoder *eos.Decoder) error {
	return m.BaseVariant.UnmarshalBinaryVariant(decoder, AtomicAttributeVariant)
}
