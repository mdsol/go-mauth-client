package go_mauth_client

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type MAuthApp struct {
	app_id          string
	rsa_private_key *rsa.PrivateKey
}

func LoadMauth(app_id string, key_file_name string) (*MAuthApp, error) {
	private_key_file, err := os.Open(key_file_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// need to convert pemfile to []byte for decoding

	pemfileinfo, _ := private_key_file.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	// read pemfile content into pembytes
	buffer := bufio.NewReader(private_key_file)
	_, err = buffer.Read(pembytes)

	// proper decoding now
	data, _ := pem.Decode([]byte(pembytes))
	privatekey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	private_key_file.Close()

	app := MAuthApp{app_id: app_id, rsa_private_key: privatekey}
	// TODO: return error?
	return &app, nil
}
