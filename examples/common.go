// This package contains two examples of use for the MAuth client
package examples

import (
	"log"
	"os"

	go_mauth_client "github.com/mdsol/go-mauth-client"
)

// Helper function to load the requisite attributes from the environment
//   MAUTH_APP_UUID - the Application UUID as provided when the application was registered
//   MAUTH_PRIVATE_KEY - the Private Key content
func LoadApp() (mauthApp *go_mauth_client.MAuthApp, err error) {
	appUUID := os.Getenv("MAUTH_APP_UUID")
	privateKeyString := os.Getenv("MAUTH_PRIVATE_KEY")

	// load the configuration from the environment
	mauthApp, err = go_mauth_client.LoadMauth(go_mauth_client.MAuthOptions{appUUID,
		privateKeyString,
		false})
	if err != nil {
		log.Fatal("Unable: to load client configuration; "+
			"did you define MAUTH_APP_UUID and "+
			"MAUTH_PRIVATE_KEY?: ", err)
	}
	return
}
