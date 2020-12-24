package api

import "encoding/json"

func unmarshallVariables(variables string) (map[string]interface{}, error) {
	var vs map[string]interface{}

	err := json.Unmarshal([]byte(variables), &vs)
	return vs, err
}
