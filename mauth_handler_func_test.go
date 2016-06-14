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
	no_json := "123"
	is_json := isJSON(no_json)
	if is_json != false {
		t.Error("Not JSON, but thinks it is")
	}
}

func TestIsJsonIsJson(t *testing.T) {
	var no_json = "{\"number\": 123}"
	is_json := isJSON(no_json)
	if is_json != true {
		t.Error("Is JSON, but thinks it is not")
	}
	no_json = "{\"number\": \"123\"}"
	is_json = isJSON(no_json)
	if is_json != true {
		t.Error("Is JSON, but thinks it is not")
	}
}
