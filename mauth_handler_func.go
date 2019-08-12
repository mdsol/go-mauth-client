package go_mauth_client

import (
	"encoding/json"
)

/*
Much of this was heavily informed by:
https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124#.9yl4dj1gd
and
https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.xj15k9f5k
*/

// GetVersion - Gets the Version for this Client
func GetVersion() string {
	return VersionString
}

// isJSON tries to work out if the content is JSON, so it can add the correct Content-Type to the Headers
// taken from http://stackoverflow.com/a/22129435/1638744
func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}
