package go_mauth_client

import "testing"

//	"net/http"
//	"net/http/httptest"
//)

/*
Influenced by:
https://elithrar.github.io/article/testing-http-handlers-go/
*/

//func TestMAuthSignerHandler(t *testing.T) {
//	req, err := http.NewRequest("GET", "https://imedidata-sandbox.imedidata.net/api/v2/users.json", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//
//
//
//}

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
