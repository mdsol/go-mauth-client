package go_mauth_client

import (
	"path/filepath"
	"testing"
	"io/ioutil"
)

func TestLoadMauth(t *testing.T) {
	app_id := "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	mauth, err := LoadMauth(app_id, filepath.Join("test", "test_mauth.priv.key"))
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
	app_id := "5ff4257e-9c16-11e0-b048-0026bbfffe5e"
	key_content, _ := ioutil.ReadFile(filepath.Join("test", "test_mauth.priv.key"))
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
