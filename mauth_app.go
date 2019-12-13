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
	AppId         string
	RsaPrivateKey *rsa.PrivateKey
	DisableV1     bool
}

type MAuthOptions struct {
	AppId      string
	PrivateKey string
	DisableV1  bool
}

// LoadMauth loads the configuration  when the private key content is in a file
func LoadMauth(options MAuthOptions) (*MAuthApp, error) {
	// Create the MAuthApp struct
	keyFileContent, err := ioutil.ReadFile(options.PrivateKey)
	if err != nil {
		keyFileContent = []byte(options.PrivateKey)
	}

	block, _ := pem.Decode(keyFileContent)
	if block == nil {
		return nil, errors.New("Unable to extract PEM content")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	app := MAuthApp{AppId: options.AppId,
		RsaPrivateKey: privateKey,
		DisableV1:     options.DisableV1}
	return &app, nil
}

// makeRequest formulates the message, including the MAuth Headers and returns a http.Request, ready to send
func (mauthApp *MAuthApp) makeRequest(method string, rawurl string, body string,
	extraHeaders map[string][]string) (req *http.Request, err error) {

	// create a new request object
	req, err = http.NewRequest(method, rawurl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	// Use the url.URL to assist with path management
	url2, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// this needs to persist
	secondsSinceEpoch := time.Now().Unix()

	if !mauthApp.DisableV1 {
		// build the MWS string
		stringToSign := MakeSignatureString(mauthApp, method, url2.Path, body, secondsSinceEpoch)

		// Sign the string
		signedString, err := SignString(mauthApp, stringToSign)
		if err != nil {
			return nil, err
		}

		// take everything and build the structure of the MAuth Headers
		madeHeaders := MakeAuthenticationHeaders(mauthApp, signedString, secondsSinceEpoch)
		for header, value := range madeHeaders {
			req.Header.Set(header, value)
		}
	}

	// build the MCC string
	stringToSignV2 := MakeSignatureStringV2(mauthApp, method, url2.Path, body, secondsSinceEpoch)

	signedStringV2, err := SignStringV2(mauthApp, stringToSignV2)
	if err != nil {
		return nil, err
	}

	madeHeadersV2 := MakeAuthenticationHeadersV2(mauthApp, signedStringV2, secondsSinceEpoch)
	for header, value := range madeHeadersV2 {
		req.Header.Set(header, value)
	}

	// Detect JSON, send appropriate Content-Type if detected
	if isJSON(body) == true {
		req.Header.Set("Content-Type", "application/json")
	}
	// Merge in any extra headers
	for header, values := range extraHeaders {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}
	// Add the User-Agent using the Client Version
	req.Header.Set("User-Agent",
		strings.Join([]string{"go-mauth-client", GetVersion()}, "/"))
	return req, nil
}
