package go_mauth_client

import (
	"crypto/rsa"
	"encoding/pem"
	"io/ioutil"
	"crypto/x509"
)

type MAuthApp struct {
	app_id          string
	rsa_private_key *rsa.PrivateKey
}

func LoadMauth(app_id string, key_file_name string) (*MAuthApp, error) {

	private_key, err := ioutil.ReadFile(key_file_name)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(private_key)

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{app_id: app_id,
		rsa_private_key: privatekey}
	// TODO: return error?
	return &app, nil
}
