package go_mauth_client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

// MAuthApp struct holds all the necessary context for a MAuth App
type MAuthApp struct {
	AppId         string
	RsaPrivateKey *rsa.PrivateKey
}

// LoadMauth loads the configuration  when the private key content is in a file
func LoadMauth(app_id string, key_file_name string) (*MAuthApp, error) {
	// Create the MAuthApp struct
	privateKey, err := ioutil.ReadFile(key_file_name)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKey)

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{AppId: app_id,
		RsaPrivateKey: privatekey}
	return &app, nil
}

// LoadMauth loads the configuration  when the private key content is passed (such as from an environment string)
func LoadMauthFromString(app_id string, key_file_content []byte) (*MAuthApp, error) {
	// Create the MAuthApp struct, when passed a byte array

	block, _ := pem.Decode(key_file_content)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{AppId: app_id,
		RsaPrivateKey: privateKey}
	return &app, nil
}
