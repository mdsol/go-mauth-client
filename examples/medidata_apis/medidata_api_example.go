// The package medidata_apis is one example of using the Go MAuth Client library to accessing a Medidata API
//
package medidata_apis

// Sample Medidata API usage

import (
	"log"

	"fmt"
	"github.com/mdsol/go-mauth-client"
	"github.com/mdsol/go-mauth-client/examples"
	"net/http"
)

// Encapsulate creation of a MAuth Client
func createClient() (mauthClient *go_mauth_client.MAuthClient, err error) {
	// load the configuration from the environment
	mauthApp, err := examples.LoadApp()
	if err != nil {
		log.Fatal("Unable: to load client configuration; "+
			"did you define MAUTH_APP_UUID and "+
			"MAUTH_PRIVATE_KEY?: ", err)
	}
	mauthClient, err = mauthApp.CreateClient("https://api.mdsol.com")
	return
}

// Get the results from the Countries API
func GetCountriesService(mauthClient *go_mauth_client.MAuthClient) (result *http.Response, err error) {
	// Access the Countries API
	targetUrl := "https://api.mdsol.com/v1/countries"

	// get the response
	result, err = mauthClient.Get(targetUrl)
	return result, err
}

func main() {
	// Create a new client
	mauthClient, err := createClient()
	if err != nil {
		log.Fatal("Error creating the client")
	}
	// Make the call and get the response
	result, err := GetCountriesService(mauthClient)
	// Handle a connection failure
	if err != nil {
		log.Fatal(fmt.Printf("Error calling the URL %s: %v", result.Request.RequestURI, err))
	}
	// Report
	log.Println("Status Code: ", result.StatusCode)
	log.Println("Results: ", result.Body)
}
