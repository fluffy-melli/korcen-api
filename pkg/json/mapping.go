// pkg/json/mapping.go

package json

import (
	"encoding/json"
)

func Mapping[T any](body []byte, t *T) (*T, error) {
	err := json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
