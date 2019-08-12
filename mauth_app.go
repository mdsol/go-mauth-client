package go_mauth_client

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MAuthApp struct holds all the necessary context for a MAuth App
type MAuthApp struct {
	AppID         string
	RsaPrivateKey *rsa.PrivateKey
}

// LoadMauth loads the configuration  when the private key content is in a file
func LoadMauth(appId string, keyFileName string) (*MAuthApp, error) {
	// Create the MAuthApp struct
	keyFileContent, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		return nil, err
	}

	// Reuse the core code
	app, err := LoadMauthFromString(appId, keyFileContent)
	return app, err
}

// LoadMauth loads the configuration  when the private key content is passed (such as from an environment string)
func LoadMauthFromString(appId string, keyFileContent []byte) (*MAuthApp, error) {
	// Create the MAuthApp struct, when passed a byte array

	block, _ := pem.Decode(keyFileContent)
	if block == nil {
		return nil, errors.New("unable to extract PEM content")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{AppID: appId,
		RsaPrivateKey: privateKey}
	return &app, nil
}

// makeRequest formulates the message, including the MAuth Headers and returns a http.Request, ready to send
func (mauthApp *MAuthApp) makeRequest(method string, rawurl string, body string,
	extraHeaders map[string][]string) (req *http.Request, err error) {
	// Use the url.URL to assist with path management
	url2, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// this needs to persist
	secondsSinceEpoch := time.Now().Unix()
	// build the MWS string
	stringToSign := MakeSignatureString(mauthApp, method, url2.Path, body, secondsSinceEpoch)
	// Sign the string
	signedString, err := SignString(mauthApp, stringToSign)
	if err != nil {
		return nil, err
	}
	// create a new request object
	req, err = http.NewRequest(method, rawurl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	// take everything and build the structure of the MAuth Headers
	madeHeaders := MakeAuthenticationHeaders(mauthApp, signedString, secondsSinceEpoch)
	for header, value := range madeHeaders {
		req.Header.Set(header, value)
	}
	// Detect JSON, send appropriate Content-Type if detected
	if isJSON(body) == true {
		req.Header.Set("Content-Type", "application/json")
	}
	// Merge in any extra headers
	for header, values := range extraHeaders {
		if len(values) == 1 {
			req.Header.Set(header, values[0])
		} else {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	}
	// Add the User-Agent using the Client Version
	req.Header.Set("User-Agent",
		strings.Join([]string{"go-mauth-client", GetVersion()}, "/"))
	return req, nil
}
