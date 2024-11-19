package utils

import "encoding/json"

func ToJSON(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
