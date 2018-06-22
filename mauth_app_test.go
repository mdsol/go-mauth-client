package go_mauth_client

import (
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
