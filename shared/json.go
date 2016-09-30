package shared

import (
	"bytes"
	"encoding/json"
)

//PrettyJSON will return the indented version of the json that was passed.
func PrettyJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
