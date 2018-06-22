// This is a simple client for the Medidata MAuth Authentication Protocol.  It can be used to access Platform Services within the Medidata Clinical Cloud.
//
// MAuth Protocol
//
// The MAuth protocol provides a fault-tolerant, service-to-service authentication scheme for Medidata and third-party applications that use web services to communicate. The Authentication Service and integrity algorithm is based on digital signatures encrypted and decrypted with a private/public key pair.
//
// The Authentication Service has two responsibilities. It provides message integrity and provenance validation by verifying a message sender's signature; its other task is to manage public keys. Each public key is associated with an application and is used to authenticate message signatures. The private key corresponding to the public key in the Authentication Service is stored by the application making a signed request; the request is encrypted with this private key. The Authentication Service has no knowledge of the application's private key, only its public key.
//
// Examples
//
// There are code examples with the methods defined in the core library.
//
// There are two code samples in the examples directory which can be used as a reference
package go_mauth_client

import (
	"net/http"
	"net/url"
	"strings"
)

// MAuthClient struct holds all the context for a MAuth Client
type MAuthClient struct {
	mauthApp *MAuthApp
	baseUrl  string
}

// CreateClient creates a MAuth Client for the baseUrl
func (mauthApp *MAuthApp) CreateClient(baseUrl string) (client *MAuthClient, err error) {
	// check for a bad baseURL
	_, err = url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}
	client = &MAuthClient{mauthApp: mauthApp, baseUrl: baseUrl}
	return
}

// fullURL returns the full URL, if we have a path it will prepend the base_url
func (mauthClient *MAuthClient) fullURL(targetUrl string) (fullUrl string, err error) {
	if strings.HasPrefix(targetUrl, "http") {
		fullUrl = targetUrl
	} else {
		// We validate the URL on create
		parsedUrl, _ := url.Parse(mauthClient.baseUrl)
		parsedUrl.Path = targetUrl
		fullUrl = parsedUrl.String()
	}
	return
}

// MAuthClient.Get executes a GET request against targetURL
func (mauthClient *MAuthClient) Get(targetURL string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("GET", fullUrl, "")
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

// MAuthClient.Delete executes a DELETE request against targetURL
func (mauthClient *MAuthClient) Delete(targetURL string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("DELETE", fullUrl, "")
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

// MAuthClient.Post executes a POST request against a targetURL
func (mauthClient *MAuthClient) Post(targetURL string, data string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("POST", fullUrl, data)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

// MAuthClient.Put executes a PUT request against a targetURL
func (mauthClient *MAuthClient) Put(targetURL string, data string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("PUT", fullUrl, data)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}
