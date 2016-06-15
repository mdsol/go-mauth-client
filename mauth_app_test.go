package main

import (
	"path/filepath"
	"io/ioutil"
	"testing"
)

func TestLoadMauth(t *testing.T) {
	mauth, err := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	if mauth.app_id != app_id {
		t.Error("App ID doesn't match")
	}
	if mauth.rsa_private_key.Validate() != nil {
		t.Error("Error validating key")
	}
}

func TestLoadMauthFromString(t *testing.T) {
	key_content, _ := ioutil.ReadFile(filepath.Join("test", "private_key.pem"))
	mauth, err := LoadMauthFromString(app_id, key_content)
	if err != nil {
		t.Error("Error creating the MAuth Struct")
	}
	if mauth.app_id != app_id {
		t.Error("App ID doesn't match")
	}
	if mauth.rsa_private_key.Validate() != nil {
		t.Error("Error validating key")
	}
}
