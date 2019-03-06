package go_mauth_client

import (
	"testing"
)

func TestIsJsonNotJson(t *testing.T) {
	content := `Platypus`
	isJson := isJSON(content)
	if isJson != false {
		t.Error("Not JSON, but thinks it is")
	}
}

func TestIsJsonIsJson(t *testing.T) {
	var content = `"123"`
	var isJson bool
	isJson = isJSON(content)
	if isJson != true {
		t.Error("Is JSON, but thinks it is not")
	}
	content = `{"number": 123}`
	isJson = isJSON(content)
	if isJson != true {
		t.Error("Is JSON, but thinks it is not")
	}
	content = `{"number": "123"}`
	isJson = isJSON(content)
	if isJson != true {
		t.Error("Is JSON, but thinks it is not")
	}
	content = `[{"number": "123"}]`
	isJson = isJSON(content)
	if isJson != true {
		t.Error("Is JSON, but thinks it is not")
	}
}
