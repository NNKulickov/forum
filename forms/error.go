package forms

import "encoding/json"

type Error struct {
	Message string `json:"message"`
}

type EmptyStruct struct {
}

type EmptyArray []struct{}

func (e EmptyArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{})
}
