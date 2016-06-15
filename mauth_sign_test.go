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

func TestMakeAuthenticationHeaders(t *testing.T) {
	const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	now := time.Now()
	secs := now.Unix()

	expected := map[string]string{
		"X-MWS-Authentication": fmt.Sprintf("MWS %s:%s", mauth_app.app_id, "some string"),
		"X-MWS-Time":           strconv.FormatInt(secs, 10),
	}
	actual := MakeAuthenticationHeaders(mauth_app, "some string", secs)
	eq := reflect.DeepEqual(expected, actual)
	if !eq {
		t.Error("Authentication headers don't match")

	}
}

func TestStringToSign(t *testing.T) {
	const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	epoch := time.Now().Unix()
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	expected := "GET" + "\n" + "/studies/123/users" + "\n" + "\n" + app_id + "\n" + strconv.FormatInt(epoch, 10)
	actual := MakeSignatureString(mauth_app, "GET", "/studies/123/users", "", -1)
	if actual != expected {
		t.Error("Signature String doesn't match")
	}
}

func TestStringToSignNoQueryParams(t *testing.T) {
	const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	epoch := time.Now().Unix()
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	expected := "GET" + "\n" + "/studies/123/users" + "\n" + "\n" + app_id + "\n" + strconv.FormatInt(epoch, 10)
	actual := MakeSignatureString(mauth_app, "GET", "/studies/123/users?until=2100", "", epoch)
	if actual != expected {
		t.Error("Signature String doesn't match: Expected ", strings.Replace(expected, "\n", " ", -1),
			"Actual ", strings.Replace(actual, "\n", " ", -1))
	}
}

func TestEpochDefinedIfMissing(t *testing.T) {
	const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	actual := MakeSignatureString(mauth_app, "GET", "/studies/123/users", "", -1)
	epoch_str := strings.Split(actual, "\n")
	epoch, _ := strconv.ParseInt(epoch_str[4], 10, 64)
	nowish := time.Unix(epoch, 0)
	now := time.Now()
	if !(now.Day() == nowish.Day() && now.Month() == nowish.Month() && now.Hour() == nowish.Hour()) {
		t.Error("Epoch not set correctly")
	}
}

func TestSignString(t *testing.T) {
	const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	const message = "Hello world"
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	actual, _ := SignString(mauth_app, message)
	expected := "IUjQhtH4C9lbCRTyca+/i4raw7ZCcyYqy5/8c79LmJcsKxkxcRuUuIdBmeUXqCDJJ25ncAs3PmRg0UzwqnQeTh5GvIqVCeRlgqZttccVhO1knbgR+sZvq2zAi5HAWycwNXNVy/r2R4/SqjTfZq4Fd/rlytBVCFLu5cigxO5yl+Gv69dgck2vNAF45jJOyS1mCbk5Zti4scy4Vca31opl9QiGiN10Z6UHXkma1fut2sGGY03Q8UDHEqNnfds1vo7NMqbeSawIVjldWhNWzbxTYM8iOocOxK5vkmj3g6Lej59pEHlnGlK/AAPUgr/soCcZoE7mYSDcDucrMj9qi4Tvlg=="
	if expected != actual {
		t.Error("Encryption does not match: ", actual)
	}
}
