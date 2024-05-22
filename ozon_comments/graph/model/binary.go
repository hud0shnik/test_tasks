package model

import "encoding/json"

func (p Post) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p Post) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (c Comment) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c Comment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}
