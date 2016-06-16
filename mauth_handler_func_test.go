package main

import "testing"

func TestIsJsonNotJson(t *testing.T) {
	content := `Platypus`
	is_json := isJSON(content)
	if is_json != false {
		t.Error("Not JSON, but thinks it is")
	}
}

func TestIsJsonIsJson(t *testing.T) {
	var content = `"123"`
	var is_json bool
	is_json = isJSON(content)
	if is_json != true {
		t.Error("Is JSON, but thinks it is not")
	}
	content = `{"number": 123}`
	is_json = isJSON(content)
	if is_json != true {
		t.Error("Is JSON, but thinks it is not")
	}
	content = `{"number": "123"}`
	is_json = isJSON(content)
	if is_json != true {
		t.Error("Is JSON, but thinks it is not")
	}
	content = `[{"number": "123"}]`
	is_json = isJSON(content)
	if is_json != true {
		t.Error("Is JSON, but thinks it is not")
	}
}
