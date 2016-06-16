package main

import (
	"path/filepath"
	"testing"
)

func TestFullURLWithRelative(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.createClient("https://innovate.mdsol.com")
	expected := "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
	// now, with a trailing slash
	client, _ = mauth_app.createClient("https://innovate.mdsol.com/")
	expected = "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ = client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen: ", actual)

	}
}

func TestFullURLWithRelativeAndParams(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.createClient("https://innovate.mdsol.com")
	expected := "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
	// now, with a trailing slash
	client, _ = mauth_app.createClient("https://innovate.mdsol.com/")
	expected = "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ = client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen: ", actual)

	}
}

func TestFullURLWithActualURL(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.createClient("https://innovate.mdsol.com")
	expected := "https://balance-innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("https://balance-innovate.mdsol.com/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
}

func TestCreateClient(t *testing.T){
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.createClient("https://innovate.mdsol.com")
	if client.base_url != "https://innovate.mdsol.com" {
		t.Error("Base URL has changed")
	}
	if client.mauth_app.app_id != app_id {
		t.Error("App ID has changed")
	}
}

