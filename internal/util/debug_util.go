package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func ConvertToMap(v any) map[string]interface{} {
	var m map[string]interface{}
	b, _ := json.Marshal(v)
	json.Unmarshal(b, &m)
	return m
}

func GetBody(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	all, _ := io.ReadAll(resp.Body)
	return string(all)
}
