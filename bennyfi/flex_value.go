// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package bennyfi

import (
	"fmt"
	"log"

	eos "github.com/eoscanada/eos-go"
)

// FlexValueVariant may hold a name, int64, asset, string, or time_point
// var FlexValueVariant = eos.NewVariantDefinition([]eos.VariantType{
// 	{Name: "monostate", Type: int64(0)},
// 	{Name: "name", Type: eos.Name("")},
// 	{Name: "string", Type: ""},
// 	{Name: "asset", Type: eos.Asset{}}, //(*eos.Asset)(nil)}, // Syntax for pointer to a type, could be any struct
// 	{Name: "time_point", Type: eos.TimePoint(0)},
// 	{Name: "microseconds", Type: Microseconds{}},
// 	{Name: "int8", Type: int8(0)},
// 	{Name: "int16", Type: int16(0)},
// 	{Name: "int32", Type: int32(0)},
// 	{Name: "int64", Type: int64(0)},
// 	{Name: "uint8", Type: uint8(0)},
// 	{Name: "uint16", Type: uint16(0)},
// 	{Name: "uint32", Type: uint32(0)},
// 	{Name: "uint64", Type: uint64(0)},
// 	{Name: "checksum256", Type: eos.Checksum256([]byte("0"))},
// })

var FlexValueVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "monostate", Type: int64(0)},
	{Name: "name", Type: eos.Name("")},
	{Name: "string", Type: ""},
	{Name: "asset", Type: eos.Asset{}}, //(*eos.Asset)(nil)}, // Syntax for pointer to a type, could be any struct
	{Name: "time_point", Type: eos.TimePoint(0)},
	{Name: "microseconds", Type: Microseconds{}},
	{Name: "int64", Type: int64(0)},
	{Name: "uint32", Type: uint32(0)},
	{Name: "uint64", Type: uint64(0)},
	{Name: "checksum256", Type: eos.Checksum256([]byte("0"))},
})

// GetVariants returns the definition of types compatible with FlexValue
func GetVariants() *eos.VariantDefinition {
	return FlexValueVariant
}

// FlexValue may hold any of the common EOSIO types
// name, int64, asset, string, time_point, or checksum256
type FlexValue struct {
	eos.BaseVariant
}

func NewFlexValue(typeId string, value interface{}) *FlexValue {
	return &FlexValue{
		BaseVariant: eos.BaseVariant{
			TypeID: GetVariants().TypeID(typeId),
			Impl:   value,
		}}
}

func FlexValueFromString(value string) *FlexValue {
	return NewFlexValue("string", value)
}
func FlexValueFromName(value eos.Name) *FlexValue {
	return NewFlexValue("name", value)
}

func FlexValueFromAsset(value eos.Asset) *FlexValue {
	return NewFlexValue("asset", value)
}

func FlexValueFromTimePoint(value eos.TimePoint) *FlexValue {
	return NewFlexValue("time_point", value)
}

func FlexValueFromMicroseconds(value Microseconds) *FlexValue {
	return NewFlexValue("microseconds", value)
}

func FlexValueFromInt64(value int64) *FlexValue {
	return NewFlexValue("int64", value)
}

func FlexValueFromUint32(value uint32) *FlexValue {
	return NewFlexValue("uint32", value)
}

func FlexValueFromUint64(value uint64) *FlexValue {
	return NewFlexValue("uint64", value)
}

func FlexValueFromChecksum(value eos.Checksum256) *FlexValue {
	return NewFlexValue("checksum256", value)
}

func (fv *FlexValue) String() string {
	switch v := fv.Impl.(type) {
	case eos.Name:
		return string(v)
	case int64, uint32, uint64:
		return fmt.Sprint(v)
	case eos.Asset:
		return v.String()
	case Microseconds:
		return v.String()
	case string:
		return v
	case eos.TimePoint:
		return v.String()
	case eos.Checksum256:
		return v.String()
	default:
		return fmt.Sprintf("received an unexpected type %T for variant %T", v, fv)
	}
}

// TimePoint returns a eos.TimePoint value of found content
func (fv *FlexValue) TimePoint() (eos.TimePoint, error) {
	switch v := fv.Impl.(type) {
	case eos.TimePoint:
		return v, nil
	default:
		return 0, &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "eos.TimePoint",
			FlexValue:    fv,
		}
	}
}

// Asset returns a string value of found content or it panics
func (fv *FlexValue) Asset() (eos.Asset, error) {
	switch v := fv.Impl.(type) {
	case eos.Asset:
		return v, nil
	default:
		return eos.Asset{}, &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "eos.Asset",
			FlexValue:    fv,
		}
	}
}

// Name returns a eos.Name value of found content or it panics
func (fv *FlexValue) Name() (eos.Name, error) {
	switch v := fv.Impl.(type) {
	case eos.Name:
		return v, nil
	case string:
		return eos.Name(v), nil
	default:
		return eos.Name(""), &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "eos.Name",
			FlexValue:    fv,
		}
	}
}

func (fv *FlexValue) Microseconds() (Microseconds, error) {
	switch v := fv.Impl.(type) {
	case Microseconds:
		return v, nil
	default:
		return Microseconds{}, &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "Microseconds",
			FlexValue:    fv,
		}
	}
}

// Int64 returns a string value of found content or it panics
func (fv *FlexValue) Int64() (int64, error) {
	switch v := fv.Impl.(type) {
	case int64:
		return v, nil
	default:
		return -1000000, &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "int64",
			FlexValue:    fv,
		}
	}
}

func (fv *FlexValue) Uint32() (uint32, error) {
	switch v := fv.Impl.(type) {
	case uint32:
		return v, nil
	default:
		return 0, &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "uint32",
			FlexValue:    fv,
		}
	}
}

func (fv *FlexValue) Uint64() (uint64, error) {
	switch v := fv.Impl.(type) {
	case uint64:
		return v, nil
	default:
		return 0, &InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for variant %T", v, fv),
			ExpectedType: "uint64",
			FlexValue:    fv,
		}
	}
}

// IsEqual evaluates if the two FlexValues have the same types and values (deep compare)
func (fv *FlexValue) IsEqual(fv2 *FlexValue) bool {

	if fv.TypeID != fv2.TypeID {
		log.Println("FlexValue types inequal: ", fv.TypeID, " vs ", fv2.TypeID)
		return false
	}

	if fv.String() != fv2.String() {
		log.Println("FlexValue Values.String() inequal: ", fv.String(), " vs ", fv2.String())
		return false
	}

	return true
}

// MarshalJSON translates to []byte
func (fv *FlexValue) MarshalJSON() ([]byte, error) {
	return fv.BaseVariant.MarshalJSON(FlexValueVariant)
}

// UnmarshalJSON translates flexValueVariant
func (fv *FlexValue) UnmarshalJSON(data []byte) error {
	return fv.BaseVariant.UnmarshalJSON(data, FlexValueVariant)
}

// UnmarshalBinary ...
func (fv *FlexValue) UnmarshalBinary(decoder *eos.Decoder) error {
	return fv.BaseVariant.UnmarshalBinaryVariant(decoder, FlexValueVariant)
}
