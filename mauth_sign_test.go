package go_mauth_client

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"

func TestMakeAuthenticationHeaders(t *testing.T) {
	mauthApp, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	now := time.Now()
	secs := now.Unix()

	expected := map[string]string{
		"X-MWS-Authentication": fmt.Sprintf("MWS %s:%s", mauthApp.AppId, "some string"),
		"X-MWS-Time":           strconv.FormatInt(secs, 10),
	}
	actual := MakeAuthenticationHeaders(mauthApp, "some string", secs)
	eq := reflect.DeepEqual(expected, actual)
	if !eq {
		t.Error("Authentication headers don't match")

	}
}

func TestStringToSign(t *testing.T) {
	epoch := time.Now().Unix()
	mauthApp, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	expected := "GET" + "\n" + "/studies/123/users" + "\n" + "\n" + app_id + "\n" + strconv.FormatInt(epoch, 10)
	actual := MakeSignatureString(mauthApp, "GET", "/studies/123/users", "", -1)
	if actual != expected {
		t.Error("Signature String doesn't match")
	}
}

func TestStringToSignNoQueryParams(t *testing.T) {
	epoch := time.Now().Unix()
	mauthApp, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	expected := "GET" + "\n" + "/studies/123/users" + "\n" + "\n" + app_id + "\n" + strconv.FormatInt(epoch, 10)
	actual := MakeSignatureString(mauthApp, "GET", "/studies/123/users?until=2100", "", epoch)
	if actual != expected {
		t.Error("Signature String doesn't match: Expected ", strings.Replace(expected, "\n", " ", -1),
			"Actual ", strings.Replace(actual, "\n", " ", -1))
	}
}

func TestEpochDefinedIfMissing(t *testing.T) {
	mauthApp, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	actual := MakeSignatureString(mauthApp, "GET", "/studies/123/users", "", -1)
	epochStr := strings.Split(actual, "\n")
	epoch, _ := strconv.ParseInt(epochStr[4], 10, 64)
	nowish := time.Unix(epoch, 0)
	now := time.Now()
	if !(now.Day() == nowish.Day() && now.Month() == nowish.Month() && now.Hour() == nowish.Hour()) {
		t.Error("Epoch not set correctly")
	}
}

func TestSignString(t *testing.T) {
	const message = "Hello world"
	mauthApp, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	actual, _ := SignString(mauthApp, message)
	expected := "ktlytWxKb6DnEaLlNRuDapLVf1DlkFCIY+/f+/VDZbH6tD6QFZ/8M3XklBEvaBqlYeACBptHZK52Yv9jbNO2gVdkQsb6Qo4467dDuHTpLmaeGTZetxZ8yuWwzHQGfqawgH9V6omnPrYbKJWaAzEdrlIxQqCWktibE6l1uW7pikZr+Y4NoDkIMavOgTdRJyOe1TDL+3GIIvIDTc5G+Mu7hNqxWRvnJTocAWFj/7ZA3GaBsHbZy9wwIzVmcloE5ahMFOlFIPI4e8DEa5sBsE7vklG25jRm8+E3GX7osslVY51RFh14KrJVIAu8gR9KzTlxRWRe8avoVf/q7CuiUyBOHA=="
	if expected != actual {
		t.Error("Encryption does not match: ", actual)
	}
}

func TestPrivateEncrypt(t *testing.T) {
	mauthApp, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	enc, err := privateEncrypt(mauthApp, []byte("encrypt_this"))
	if err != nil {
		t.Error(err)
	}
	if len(enc) != 256 {
		t.Error("Wrong size of encrypted data")
	}
}
