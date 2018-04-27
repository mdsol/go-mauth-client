package examples

import (
	"github.com/mdsol/go-mauth-client"
	"log"
	"fmt"
)

/*
A Sample API Gateway usage
*/


 func createClient()(mauthClient *go_mauth_client.MAuthClient, err error){
 	// load the configuration from the environment
 	mauthApp, err := loadApp()
 	if err != nil {
		log.Fatal("Unable: to load client configuration; " +
			"did you define MAUTH_APP_UUID and " +
			"MAUTH_PRIVATE_KEY?: ", err)
	}
	mauthClient, err = mauthApp.CreateClient("https://apigw.imedidata.com")
 	return
 }

func main()  {
	targetUrl := "https://apigw.imedidata.com/v1/countries"

	mauthClient, err := createClient()
	if err != nil {
		log.Fatal("Error creating the client")
	}
	result, err := mauthClient.Get(targetUrl)
	if err != nil {
		log.Fatal(fmt.Printf("Error calling the URL %s: %v", targetUrl, err))
	}
	log.Println("Status Code: ", result.StatusCode)
	log.Println("Results: ", result.Body)
}
