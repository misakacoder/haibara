package util

import "encoding/json"

func ToJSONString(object any) (string, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func ParseObject[T any](text string) (T, error) {
	var object T
	err := json.Unmarshal([]byte(text), &object)
	return object, err
}
