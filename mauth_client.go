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

	"github.com/pkg/errors"
)

// MAuthClient struct holds all the context for a MAuth Client
type MAuthClient struct {
	mauthApp     *MAuthApp
	baseURL      *url.URL
	extraHeaders map[string][]string
}

// CreateClient creates a MAuth Client for the baseUrl
func (mauthApp *MAuthApp) CreateClient(baseURL string) (client *MAuthClient, err error) {
	// check for a bad baseURL
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}
	client = &MAuthClient{mauthApp: mauthApp,
		baseURL:      parsedURL,
		extraHeaders: make(map[string][]string)}
	return
}

// SetHeader - adds a Header to the Request
func (mauthClient *MAuthClient) SetHeader(headerName, headerValue string) {
	header, exists := mauthClient.extraHeaders[headerName]
	if !exists {
		header = []string{}
	}
	header = append(header, headerValue)
	mauthClient.extraHeaders[headerName] = header
}

// GetHeader - get a Header by Name
func (mauthClient *MAuthClient) GetHeader(headerName string) ([]string, error) {
	header, exists := mauthClient.extraHeaders[headerName]
	if !exists {
		return nil, errors.New("No such header " + headerName)
	}
	return header, nil
}

// GetHeaders - get Headers
func (mauthClient *MAuthClient) GetHeaders() map[string][]string {
	return mauthClient.extraHeaders
}

// fullURL returns the full URL, if we have a path it will prepend the base_url
func (mauthClient *MAuthClient) fullURL(targetURL string) (fullURL string, err error) {
	var parsedURL *url.URL
	if strings.HasPrefix(targetURL, "http") {
		// an entire URL
		parsedURL, err = url.Parse(targetURL)
		if err != nil {
			return "", err
		}
	} else {
		parsedURL = mauthClient.baseURL
		// a partial URL
		parsedURL.Path = targetURL
	}
	fullURL = parsedURL.String()
	return
}

// Get executes a GET request against targetURL
func (mauthClient *MAuthClient) Get(targetURL string) (response *http.Response, err error) {
	fullURL, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("GET", fullURL, "", mauthClient.extraHeaders)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	response, err = client.Do(req)
	return
}

// Delete executes a DELETE request against targetURL
func (mauthClient *MAuthClient) Delete(targetURL string) (response *http.Response, err error) {
	fullURL, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("DELETE", fullURL, "", mauthClient.extraHeaders)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

// Post executes a POST request against a targetURL
func (mauthClient *MAuthClient) Post(targetURL string, data string) (response *http.Response, err error) {
	fullURL, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("POST", fullURL, data, mauthClient.extraHeaders)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

// Put executes a PUT request against a targetURL
func (mauthClient *MAuthClient) Put(targetURL string, data string) (response *http.Response, err error) {
	fullURL, err := mauthClient.fullURL(targetURL)
	if err != nil {
		return nil, err
	}
	req, err := mauthClient.mauthApp.makeRequest("PUT", fullURL, data, mauthClient.extraHeaders)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}
