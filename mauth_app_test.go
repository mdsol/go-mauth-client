package go_mauth_client

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestLoadMauth(t *testing.T) {
	mauth, err := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	if mauth.AppId != app_id {
		t.Error("App ID doesn't match")
	}
	if mauth.RsaPrivateKey.Validate() != nil {
		t.Error("Error validating key")
	}
}

func TestLoadMauthMissingFile(t *testing.T) {
	_, err := LoadMauth(app_id, filepath.Join("test", "banana.pem"))
	if err == nil {
		t.Error("Expected Error creating the MAuth Struct")
	}
}

func TestLoadMauthNotKey(t *testing.T) {
	_, err := LoadMauth(app_id, filepath.Join("test", "junk.pem"))
	if err == nil {
		t.Error("Expected Error loading an empty Private Key file")
	}
}

func TestLoadMauthInvalidKey(t *testing.T) {
	_, err := LoadMauth(app_id, filepath.Join("test", "invalid.pem"))
	if err == nil {
		t.Error("Expected Error loading an invalid Private Key")
	}
}

func TestLoadMauthFromString(t *testing.T) {
	keyContent, _ := ioutil.ReadFile(filepath.Join("test", "private_key.pem"))
	mauth, err := LoadMauthFromString(app_id, keyContent)
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	if mauth.AppId != app_id {
		t.Error("App ID doesn't match")
	}
	if mauth.RsaPrivateKey.Validate() != nil {
		t.Error("Error validating key")
	}
}

func TestLoadMauthFromStringNotKey(t *testing.T) {
	keyContent, _ := ioutil.ReadFile(filepath.Join("test", "junk.pem"))
	_, err := LoadMauthFromString(app_id, keyContent)
	if err == nil {
		t.Error("Expected Error loading an empty Private Key file")
	}
}

func TestLoadMauthFromStringInvalidKey(t *testing.T) {
	keyContent, _ := ioutil.ReadFile(filepath.Join("test", "invalid.pem"))
	_, err := LoadMauthFromString(app_id, keyContent)
	if err == nil {
		t.Error("Expected Error loading an invalid Private Key")
	}
}

// Example of loading the MAuth configuration from a path
func ExampleLoadMauth() {
	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth client
	var client *MAuthApp
	client, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create client: ", err)
	}
	println("Created MAuth App for APP_UUID ", client.AppId)
}

// Example for creating client where the Private Key is provided as a String
func ExampleLoadMauthFromString() {
	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and the content of the Private Key
	var keyString []byte
	keyString, _ = ioutil.ReadFile(filepath.Join("test", "private_key.pem"))
	// create a MAuth client
	var client *MAuthApp
	client, err := LoadMauthFromString(appUUID, keyString)
	if err != nil {
		log.Fatal("Unable to create client: ", err)
	}
	println("Created MAuth App for APP_UUID ", client.AppId)
}

func TestMAuthApp_makeRequestUserAgent(t *testing.T) {
	mauth, err := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	expected := fmt.Sprintf("go-mauth-client/%s", GetVersion())
	extraHeaders := make(map[string][]string)
	request, err := mauth.makeRequest("GET", "/some/url", "", extraHeaders)
	if expected != request.Header.Get("User-Agent") {
		t.Error("Expected User-Agent to be", expected, "got",
			request.Header.Get("User-Agent"))
	}

}

func TestMAuthApp_makeRequestEmpty(t *testing.T) {
	mauth, err := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	extraHeaders := make(map[string][]string)
	request, err := mauth.makeRequest("GET", "/some/url", "", extraHeaders)
	if "" != request.Header.Get("Content-Type") {
		t.Error("Expected Content-type to be empty got ",
			request.Header.Get("Content-Type"))
	}

}

func TestMAuthApp_makeRequestJSON(t *testing.T) {
	mauth, err := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	extraHeaders := make(map[string][]string)
	request, err := mauth.makeRequest("POST", "/some/url", `{"app_id":12345}`, extraHeaders)
	if "application/json" != request.Header.Get("Content-type") {
		t.Error("Expected Content-type to be 'application/json' got '",
			request.Header.Get("Content-Type"), "'")
	}
}

func TestMAuthApp_makeRequestInvalidURL(t *testing.T) {
	mauth, err := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	extraHeaders := make(map[string][]string)
	_, err = mauth.makeRequest("POST", "\x7f\x8c\x98", `{"app_id":12345}`, extraHeaders)
	if err == nil {
		t.Error("Expected Error with non-sensical URL")
	}
}
