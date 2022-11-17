package testutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func JsonDecode(t *testing.T, input string, obj any) {
	data := []byte(input)
	if err := json.Unmarshal(data, &obj); err != nil {
		assert.Fail(t, "could not unmarshal json: %v", err)
	}
}

func JsonEncode(t *testing.T, obj any) string {
	b, err := json.Marshal(obj)
	if err != nil {
		assert.Fail(t, "could not marshal json %v", err)
	}
	return string(b)
}
