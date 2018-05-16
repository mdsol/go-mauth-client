package examples

import (
	"log"
	"os"

	"github.com/mdsol/go-mauth-client"
)

func LoadApp() (mauthApp *go_mauth_client.MAuthApp, err error) {
	appUUID := os.Getenv("MAUTH_APP_UUID")
	privateKeyString := os.Getenv("MAUTH_PRIVATE_KEY")

	// load the configuration from the environment
	mauthApp, err = go_mauth_client.LoadMauthFromString(appUUID,
		[]byte(privateKeyString))
	if err != nil {
		log.Fatal("Unable: to load client configuration; "+
			"did you define MAUTH_APP_UUID and "+
			"MAUTH_PRIVATE_KEY?: ", err)
	}
	return
}
