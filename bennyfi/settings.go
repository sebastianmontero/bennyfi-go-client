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
	"strconv"

	eos "github.com/eoscanada/eos-go"
)

var (
	SettingVRFContract = "VRF_CONTRACT"
)

type Setting struct {
	ID          uint64             `json:"id"`
	Key         string             `json:"key"`
	Values      []FlexValue        `json:"values"`
	CreatedDate eos.BlockTimestamp `json:"created_date"`
	UpdatedDate eos.BlockTimestamp `json:"updated_date"`
}

func (m *BennyfiContract) setter(owner eos.AccountName,
	key string, flexValue *FlexValue, action eos.ActionName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["setter"] = owner
	actionData["key"] = key
	actionData["value"] = flexValue
	_, err := m.Contract.ExecAction(string(owner), string(action), actionData)
	if err != nil {
		return "", err
	}
	return "", nil
	// return m.ExecAction(owner, string(action), actionData)
}

func (m *BennyfiContract) SetSetting(owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.setter(owner, key, flexValue, eos.ActN("setsetting"))
}

func (m *BennyfiContract) AppendSetting(owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.setter(owner, key, flexValue, eos.ActN("appndsetting"))
}

func (m *BennyfiContract) ClipSetting(owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.setter(owner, key, flexValue, eos.ActN("clipsetting"))
}

func (m *BennyfiContract) EraseSetting(owner eos.AccountName, key string) (string, error) {

	actionData := make(map[string]interface{})
	actionData["setter"] = owner
	actionData["key"] = key

	return m.ExecAction(owner, "erasesetting", actionData)
}

func (m *BennyfiContract) GetSettings() ([]Setting, error) {
	return m.GetSettingsReq(nil)
}

func (m *BennyfiContract) GetSettingsReq(req *eos.GetTableRowsRequest) ([]Setting, error) {

	var settings []Setting
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "settings"
	err := m.GetTableRows(*req, &settings)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return settings, nil
}

func (m *BennyfiContract) SetupConfigSettings(owner eos.AccountName, settings interface{}) error {
	for _, value := range settings.([]interface{}) {
		err := m.SetupConfigSetting(owner, value.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *BennyfiContract) SetupConfigSetting(owner eos.AccountName, setting map[interface{}]interface{}) error {

	fv, err := StringToSetting(setting["type"].(string), setting["value"].(string))
	if err != nil {
		return err
	}
	_, err = m.SetSetting(owner, setting["key"].(string), &fv)
	if err != nil {
		return err
	}

	return nil
}

func StringToSetting(settingType, stringValue string) (FlexValue, error) {
	if settingType == "string" {
		return FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: GetVariants().TypeID("string"),
				Impl:   stringValue,
			},
		}, nil
	} else if settingType == "name" {
		return FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: GetVariants().TypeID("name"),
				Impl:   eos.Name(stringValue),
			},
		}, nil
	} else if settingType == "int64" {
		i, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return FlexValue{}, fmt.Errorf("cannot convert settings value to int64: %v", err)
		}
		return FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: GetVariants().TypeID("int64"),
				Impl:   int64(i),
			},
		}, nil
	} else if settingType == "asset" {
		a, err := eos.NewAssetFromString(stringValue)
		if err != nil {
			return FlexValue{}, fmt.Errorf("cannot convert settings value to asset: %v", err)
		}
		return FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: GetVariants().TypeID("asset"),
				Impl:   a,
			},
		}, nil
	} else if settingType == "timepoint" {
		fmt.Println("type: timepoint is not yet supported from CLI")
		return FlexValue{}, fmt.Errorf("not yet supported")
	}
	return FlexValue{}, fmt.Errorf("unsupported settings data type: %v", settingType)
}

// InvalidTypeError is used the type of a FlexValue doesn't match expectations
type InvalidTypeError struct {
	Label        string
	ExpectedType string
	FlexValue    *FlexValue
}

func (c *InvalidTypeError) Error() string {
	return fmt.Sprintf("received an unexpected type %T for metadata variant %T", c.ExpectedType, c.FlexValue)
}

// FlexValueVariant may hold a name, int64, asset, string, or time_point
var FlexValueVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "monostate", Type: int64(0)},
	{Name: "name", Type: eos.Name("")},
	{Name: "string", Type: ""},
	{Name: "asset", Type: eos.Asset{}}, //(*eos.Asset)(nil)}, // Syntax for pointer to a type, could be any struct
	{Name: "time_point", Type: eos.TimePoint(0)},
	{Name: "int64", Type: int64(0)},
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

func (fv *FlexValue) String() string {
	switch v := fv.Impl.(type) {
	case eos.Name:
		return string(v)
	case int64:
		return fmt.Sprint(v)
	case eos.Asset:
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

// Name returns a string value of found content or it panics
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
