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
	"strconv"
	"time"

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

func (m *Setting) ValueCount() int {
	return len(m.Values)
}

func (m *Setting) Get(pos int) (FlexValue, error) {
	if pos < m.ValueCount() {
		return m.Values[0], nil
	}
	return FlexValue{}, fmt.Errorf("setting has no value at pos: %v, value count: %v", pos, m.ValueCount())
}

func (m *Setting) GetAsName(pos int) (eos.Name, error) {
	v, err := m.Get(pos)
	if err != nil {
		return eos.Name(""), err
	}
	return v.Name(), nil
}

func (m *Setting) GetAsString(pos int) (string, error) {
	v, err := m.Get(pos)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

func (m *Setting) GetAsTimePoint(pos int) (eos.TimePoint, error) {
	v, err := m.Get(pos)
	if err != nil {
		return eos.TimePoint(0), err
	}
	return v.TimePoint(), nil
}

func (m *Setting) GetAsAsset(pos int) (eos.Asset, error) {
	v, err := m.Get(pos)
	if err != nil {
		return eos.Asset{}, err
	}
	return v.Asset(), nil
}

func (m *Setting) GetAsInt64(pos int) (int64, error) {
	v, err := m.Get(pos)
	if err != nil {
		return 0, err
	}
	return v.Int64(), nil
}

func (m *BennyfiContract) setter(owner eos.AccountName,
	key string, flexValue *FlexValue, action eos.ActionName) (string, error) {
	actionData := m.getSetterData(owner, key, flexValue)
	return m.ExecAction(string(owner), string(action), actionData)
}

func (m *BennyfiContract) proposeSetter(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName,
	key string, flexValue *FlexValue, action eos.ActionName) (string, error) {
	actionData := m.getSetterData(owner, key, flexValue)
	return m.ProposeAction(proposerName, requested, expireIn, string(owner), string(action), actionData)
}

func (m *BennyfiContract) getSetterData(owner eos.AccountName,
	key string, flexValue *FlexValue) map[string]interface{} {
	actionData := make(map[string]interface{})
	actionData["setter"] = owner
	actionData["key"] = key
	actionData["value"] = flexValue
	return actionData
}

func (m *BennyfiContract) SetSetting(owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.setter(owner, key, flexValue, eos.ActN("setsetting"))
}

func (m *BennyfiContract) ProposeSetSetting(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.proposeSetter(proposerName, requested, expireIn, owner, key, flexValue, eos.ActN("setsetting"))
}

func (m *BennyfiContract) AppendSetting(owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.setter(owner, key, flexValue, eos.ActN("appndsetting"))
}

func (m *BennyfiContract) ProposeAppendSetting(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.proposeSetter(proposerName, requested, expireIn, owner, key, flexValue, eos.ActN("appndsetting"))
}

func (m *BennyfiContract) ClipSetting(owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.setter(owner, key, flexValue, eos.ActN("clipsetting"))
}

func (m *BennyfiContract) ProposeClipSetting(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName,
	key string, flexValue *FlexValue) (string, error) {

	return m.proposeSetter(proposerName, requested, expireIn, owner, key, flexValue, eos.ActN("clipsetting"))
}

func (m *BennyfiContract) EraseSetting(owner eos.AccountName, key string) (string, error) {

	actionData := make(map[string]interface{})
	actionData["setter"] = owner
	actionData["key"] = key

	return m.ExecAction(owner, "erasesetting", actionData)
}

func (m *BennyfiContract) ProposeEraseSetting(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName,
	key string) (string, error) {

	actionData := make(map[string]interface{})
	actionData["setter"] = owner
	actionData["key"] = key
	return m.ProposeAction(proposerName, requested, expireIn, string(owner), "erasesetting", actionData)
}

func (m *BennyfiContract) GetSettings() ([]Setting, error) {
	return m.GetSettingsReq(nil)
}

func (m *BennyfiContract) GetSetting(key string) (*Setting, error) {
	settings, err := m.GetSettings()
	if err != nil {
		return nil, err
	}
	for _, setting := range settings {
		if setting.Key == key {
			return &setting, nil
		}
	}
	return nil, nil
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

func (m *BennyfiContract) SetupConfigSetting(owner eos.AccountName, configSetting map[interface{}]interface{}) error {

	setting, err := GetConfigSetting(configSetting)
	if err != nil {
		return err
	}
	_, err = m.SetSetting(owner, setting.Key, &setting.Values[0])
	if err != nil {
		return err
	}

	return nil
}

func (m *BennyfiContract) ProposeConfigSettings(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName, settings interface{}) error {
	for _, value := range settings.([]interface{}) {
		err := m.ProposeConfigSetting(proposerName, requested, expireIn, owner, value.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *BennyfiContract) ProposeConfigSetting(proposerName interface{}, requested []eos.PermissionLevel, expireIn time.Duration, owner eos.AccountName, configSetting map[interface{}]interface{}) error {

	setting, err := GetConfigSetting(configSetting)
	if err != nil {
		return err
	}
	_, err = m.ProposeSetSetting(proposerName, requested, expireIn, owner, setting.Key, &setting.Values[0])
	if err != nil {
		return err
	}

	return nil
}

func GetConfigSetting(setting map[interface{}]interface{}) (*Setting, error) {

	fv, err := StringToSetting(setting["type"].(string), setting["value"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed parsing string setting to flex value, error: %v", err)
	}
	return &Setting{
		Key: setting["key"].(string),
		Values: []FlexValue{
			fv,
		},
	}, nil
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
	} else if settingType == "uint32" {
		i, err := strconv.ParseUint(stringValue, 10, 32)
		if err != nil {
			return FlexValue{}, fmt.Errorf("cannot convert settings value to uint32: %v", err)
		}
		return FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: GetVariants().TypeID("uint32"),
				Impl:   i,
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
