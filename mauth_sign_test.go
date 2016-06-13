package go_mauth_client

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
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
		"X-MWS-Time":           string(secs),
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
	actual := MakeSignatureString(mauth_app, "GET", "/studies/123/users", "")
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
		t.Error("Signature String doesn't match")
	}
}

func TestEpochDefinedIfMissing(t *testing.T) {
	const app_id = "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	actual := MakeSignatureString(mauth_app, "GET", "/studies/123/users", "")
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
	const message = "string to sign"
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	signature, _ := SignString(mauth_app, message)
	hashed := sha512.Sum512([]byte(message))
	err := rsa.VerifyPKCS1v15(&mauth_app.rsa_private_key.PublicKey, crypto.SHA512, hashed[:], signature)
	if err != nil {
		t.Error("Error verifying signature")
	}
}
