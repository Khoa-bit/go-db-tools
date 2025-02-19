package tool

import (
	"encoding/json"
	"log"
)

// DebugMarshal returns a pretty printed JSON string of the object.
func DebugMarshal(obj interface{}) string {
	objJSON, err := json.MarshalIndent(obj, " ", "  ")
	if err != nil {
		log.Fatalf(err.Error(), obj)
	}

	return string(objJSON)
}
