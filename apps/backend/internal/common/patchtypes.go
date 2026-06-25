package common

import "encoding/json"

type PatchStringField struct {
	Value string
	IsSet bool
	Null  bool
}

func (f *PatchStringField) UnmarshalJSON(data []byte) error {
	f.IsSet = true
	if string(data) == "null" {
		f.Null = true
		f.Value = ""

		return nil
	}

	return json.Unmarshal(data, &f.Value)
}

type PatchIntField struct {
	Value int
	IsSet bool
	Null  bool
}

func (f *PatchIntField) UnmarshalJSON(data []byte) error {
	f.IsSet = true
	if string(data) == "null" {
		f.Null = true
		f.Value = 0

		return nil
	}

	return json.Unmarshal(data, &f.Value)
}

type PatchStringSliceField struct {
	Value []string
	IsSet bool
	Null  bool
}

func (f *PatchStringSliceField) UnmarshalJSON(data []byte) error {
	f.IsSet = true
	if string(data) == "null" {
		f.Null = true
		f.Value = nil

		return nil
	}

	return json.Unmarshal(data, &f.Value)
}

type PatchNullableStringField struct {
	Value *string
	IsSet bool
}

func (f *PatchNullableStringField) UnmarshalJSON(data []byte) error {
	f.IsSet = true
	if string(data) == "null" {
		f.Value = nil

		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	f.Value = &value

	return nil
}
