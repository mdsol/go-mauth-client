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

// CreateClient creates a MAuth Client
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

// MAuthClient.Get executes a GET request against a URL
func (mauthClient *MAuthClient) Get(target_url string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(target_url)
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

// MAuthClient.Delete executes a DELETE request against a URL
func (mauthClient *MAuthClient) Delete(target_url string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(target_url)
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

// MAuthClient.Post executes a POST request against a URL
func (mauthClient *MAuthClient) Post(target_url string, data string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(target_url)
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

// MAuthClient.Put executes a PUT request against a URL
func (mauthClient *MAuthClient) Put(target_url string, data string) (response *http.Response, err error) {
	fullUrl, err := mauthClient.fullURL(target_url)
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
