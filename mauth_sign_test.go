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
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
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

func TestMakeAuthenticationHeadersV2(t *testing.T) {
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	now := time.Now()
	secs := now.Unix()

	expected := map[string]string{
		"MCC-Authentication": fmt.Sprintf("MWSV2 %s:%s;", mauthApp.AppId, "some string"),
		"MCC-Time":           strconv.FormatInt(secs, 10),
	}
	actual := MakeAuthenticationHeadersV2(mauthApp, "some string", secs)
	eq := reflect.DeepEqual(expected, actual)
	if !eq {
		t.Error("Authentication headers don't match")

	}
}

func TestStringToSign(t *testing.T) {
	epoch := time.Now().Unix()
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	expected := "GET" + "\n" + "/studies/123/users" + "\n" + "\n" + app_id + "\n" + strconv.FormatInt(epoch, 10)
	actual := MakeSignatureString(mauthApp, "GET", "/studies/123/users", "", -1)
	if actual != expected {
		t.Error("Signature String doesn't match")
	}
}

func TestStringToSignV2(t *testing.T) {
	epoch := time.Now().Unix()
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	expected := "GET" + "\n" + "/studies/123/users" + "\n" +
		"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e" + "\n" +
		app_id + "\n" + strconv.FormatInt(epoch, 10) + "\n"
	actual := MakeSignatureStringV2(mauthApp, "GET", "/studies/123/users", "", -1)
	if actual != expected {
		t.Error("Signature String doesn't match")
	}
}

func TestStringToSignNoQueryParams(t *testing.T) {
	epoch := time.Now().Unix()
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	expected := "GET" + "\n" + "/studies/123/users" + "\n" + "\n" + app_id + "\n" + strconv.FormatInt(epoch, 10)
	actual := MakeSignatureString(mauthApp, "GET", "/studies/123/users?until=2100", "", epoch)
	if actual != expected {
		t.Error("Signature String doesn't match: Expected ", strings.Replace(expected, "\n", " ", -1),
			"Actual ", strings.Replace(actual, "\n", " ", -1))
	}
}

func TestStringToSignV2QueryParams(t *testing.T) {
	epoch := time.Now().Unix()
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	expected := "GET" + "\n" + "/studies/123/users" + "\n" +
		"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e" + "\n" +
		app_id + "\n" + strconv.FormatInt(epoch, 10) + "\n" + "foo=bar1&foo=bar2&until=2100"
	actual := MakeSignatureStringV2(mauthApp, "GET", "/studies/123/users?until=2100&foo=bar2&foo=bar1", "", epoch)
	if actual != expected {
		t.Error("Signature String doesn't match: Expected ", strings.Replace(expected, "\n", " ", -1),
			"Actual ", strings.Replace(actual, "\n", " ", -1))
	}
}

func TestEpochDefinedIfMissing(t *testing.T) {
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	actual := MakeSignatureString(mauthApp, "GET", "/studies/123/users", "", -1)
	epochStr := strings.Split(actual, "\n")
	epoch, _ := strconv.ParseInt(epochStr[4], 10, 64)
	nowish := time.Unix(epoch, 0)
	now := time.Now()
	if !(now.Day() == nowish.Day() && now.Month() == nowish.Month() && now.Hour() == nowish.Hour()) {
		t.Error("Epoch not set correctly")
	}
}

func TestEpochDefinedIfMissingV2(t *testing.T) {
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	actual := MakeSignatureStringV2(mauthApp, "GET", "/studies/123/users", "", -1)
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
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	actual, _ := SignString(mauthApp, message)
	expected := "ktlytWxKb6DnEaLlNRuDapLVf1DlkFCIY+/f+/VDZbH6tD6QFZ/8M3XklBEvaBqlYeACBptHZK52Yv9jbNO2gVdkQsb6Qo4467dDuHTpLmaeGTZetxZ8yuWwzHQGfqawgH9V6omnPrYbKJWaAzEdrlIxQqCWktibE6l1uW7pikZr+Y4NoDkIMavOgTdRJyOe1TDL+3GIIvIDTc5G+Mu7hNqxWRvnJTocAWFj/7ZA3GaBsHbZy9wwIzVmcloE5ahMFOlFIPI4e8DEa5sBsE7vklG25jRm8+E3GX7osslVY51RFh14KrJVIAu8gR9KzTlxRWRe8avoVf/q7CuiUyBOHA=="
	if expected != actual {
		t.Error("Encryption does not match: ", actual)
	}
}

func TestSignStringV2(t *testing.T) {
	const message = "Hello world"
	mauthApp, _ := LoadMauth(MAuthOptions{app_id, filepath.Join("test", "private_key.pem"), false})
	actual, _ := SignStringV2(mauthApp, message)
	expected := "KIFAy3SScS3gNX9m4cRQpyR+BoOKQAVQszWawwlN4ad/x4HUyHVXiNwIiSV+PecXsclexIzLlmanUqzSZgaJBR8u9gXlKi+XMYPFT7zceSl6q2hokiuZAg98tCsjvBExTKOKwfgbsuK033xPyiqDD64a8m2xAEw5eU98NfJ2E0m6WMD/ACzua8l66KgxdusGKmi6ZyPkc0Z+qpX3ig2Xu1eooqaUyAA+OqyJIziepgbMSXTiGoVyxhLEwfNy13BS774r6L0QEKcPwnIJ3Zu1D8ZslhSQv6V57YYPt//BfBl0B9gylfb8i+PQtHTXoXHM5bMrlkzTcqjWO2epeIENPg=="
	if expected != actual {
		t.Error("Encryption does not match: ", actual)
	}
}
