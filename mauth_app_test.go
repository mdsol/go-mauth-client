package go_mauth_client

import (
	"io/ioutil"
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
