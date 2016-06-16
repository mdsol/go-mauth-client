package main

import (
	"net/http"
	"strings"
	"net/url"
)


type MAuthClient struct {
	mauth_app *MAuthApp
	base_url string
}

// create a client
func (mauth_app *MAuthApp) createClient(base_url string)(client MAuthClient, err error){
	client = MAuthClient{mauth_app: mauth_app, base_url: base_url}
	return
}

// fullURL returns the full URL, if we have a path it will prepend the base_url
func (mauth_client *MAuthClient) fullURL(target_url string)(full_url string, err error){
	if strings.HasPrefix(target_url, "http") {
		full_url = target_url
	} else {
		parsed_url, err := url.Parse(mauth_client.base_url)
		if err != nil {
			return "", err
		}
		parsed_url.Path = target_url
		full_url = parsed_url.String()
	}
	return
}

// MAuthClient.get executes a GET request against a URL
func (mauth_client *MAuthClient) get(target_url string)(response *http.Response, err error){
	req, err := mauth_client.mauth_app.makeRequest("GET", target_url, "")
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

// MAuthClient.post executes a POST request against a URL
func (mauth_client *MAuthClient) post(target_url string, data string)(response *http.Response, err error){
	req, err := mauth_client.mauth_app.makeRequest("POST", target_url, data)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	response, err = client.Do(req)
	return
}

