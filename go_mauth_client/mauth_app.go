package go_mauth_client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

type MAuthApp struct {
	App_ID          string
	RSA_Private_Key *rsa.PrivateKey
}

func LoadMauth(app_id string, key_file_name string) (*MAuthApp, error) {
	// Create the MAuthApp struct
	private_key, err := ioutil.ReadFile(key_file_name)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(private_key)

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{App_ID: app_id,
		RSA_Private_Key: privatekey}
	return &app, nil
}

func LoadMauthFromString(app_id string, key_file_content []byte) (*MAuthApp, error) {
	// Create the MAuthApp struct, when passed a byte array

	block, _ := pem.Decode(key_file_content)

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{App_ID: app_id,
		RSA_Private_Key: privatekey}
	return &app, nil
}
