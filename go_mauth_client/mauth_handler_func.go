package go_mauth_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

/*
Much of this was heavily informed by:
https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124#.9yl4dj1gd
and
https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.xj15k9f5k
*/

// taken from http://stackoverflow.com/a/22129435/1638744
func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func (mauth_app *MAuthApp) makeRequest(method string, rawurl string, body string) (req *http.Request, err error) {
	// Use the url.URL to assist with path management
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// this needs to persist
	seconds_since_epoch := time.Now().Unix()
	// build the MWS string
	string_to_sign := MakeSignatureString(mauth_app, method, url.Path, body, seconds_since_epoch)
	// Sign the string
	signed_string, err := SignString(mauth_app, string_to_sign)
	if err != nil {
		return nil, err
	}
	// create a new request object
	req, err = http.NewRequest(method, rawurl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	// take everything and build the structure of the MAuth Headers
	made_headers := MakeAuthenticationHeaders(mauth_app, signed_string, seconds_since_epoch)
	for header, value := range made_headers {
		req.Header.Set(header, value)
	}
	// Detect JSON, send appropriate Content-Type if detected
	if isJSON(body) {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}
