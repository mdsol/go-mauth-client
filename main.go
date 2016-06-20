package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mdsol/go-mauth-client/go_mauth_client"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type ApplicationContext struct {
	app_uuid         string
	private_key_text []byte
	private_key_file string
}

func CheckAction(action *string) bool {
	switch *action {
	case "GET", "POST", "PUT", "DELETE":
		return true
	default:
		return false
	}
}

// Process the Configuration Content (Easier to test)
func ProcessConfiguration(content []byte) (mauth_app *go_mauth_client.MAuthApp, err error) {
	var context map[string]string
	err = json.Unmarshal(content, &context)
	if err != nil {
		return nil, err
	}
	app_uuid := context["app_uuid"]
	private_key_file := context["private_key_file"]
	private_key_text := context["private_key_text"]
	if IsNull(&app_uuid) {
		return nil, errors.New("Need an app_uuid specified")
	}

	// TODO: check for private_key_text
	if IsNull(&private_key_file) && IsNull(&private_key_text) {
		return nil, errors.New("Need a key specified")
	}
	if IsNull(&private_key_file) {
		// read from the embedded value
		mauth_app, err = go_mauth_client.LoadMauthFromString(app_uuid,
			[]byte(private_key_text))
	} else {
		// load the key from a file
		mauth_app, err = go_mauth_client.LoadMauth(app_uuid,
			private_key_file)
	}
	return
}

// Load the Configuration from a JSON file
// why JSON, you may ask, using stdlib as much as possible
func LoadMAuthConfig(file_name string) (mauth_app *go_mauth_client.MAuthApp, err error) {
	if _, err = os.Stat(file_name); os.IsNotExist(err) {
		return nil, err
	}
	content, err := ioutil.ReadFile(file_name)
	if err != nil {
		return nil, err
	}
	mauth_app, err = ProcessConfiguration(content)
	return
}

func IsNull(value *string) bool {
	return *value == ""
}

func main() {
	config_file := flag.String("config", "", "Specify the configuration file")

	key_file := flag.String("private-key", "", "Specify the private key file")
	app_uuid := flag.String("app-uuid", "", "Specify the App UUID")

	action := flag.String("method", "GET", "Specify the method (GET, POST, PUT, DELETE)")
	data := flag.String("data", "", "Specify the data")
	headers := flag.Bool("headers", false, "Print the Response Headers")

	flag.Parse()
	// No information supplied
	if IsNull(config_file) && (IsNull(key_file) || IsNull(app_uuid)) {
		println("Need to specify configuration file or app settings")
		flag.Usage()
		os.Exit(1)
	}
	var mauth_app *go_mauth_client.MAuthApp
	// Load the MAuth Config
	if !IsNull(config_file) {
		// a config file has been passed
		var err error
		mauth_app, err = LoadMAuthConfig(*config_file)
		if err != nil {
			println("Error loading configuration: ", err)
			os.Exit(1)
		}
	} else {
		var err error
		mauth_app, err = go_mauth_client.LoadMauth(*app_uuid, *key_file)
		if err != nil {
			println("Error loading configuration: ", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Created MAuth App with App UUID: %s\n", mauth_app.App_ID)
	action_matches := CheckAction(action)
	if !action_matches {
		println("Action ", action, "is not known")
		flag.Usage()
		os.Exit(1)
	}
	var args []string
	args = flag.Args()
	if len(args) == 0 {
		println("No URL, nothing to do")
		flag.Usage()
		os.Exit(1)
	}
	target_url, err := url.Parse(args[0])
	if err != nil {
		println("Unable to parse url: ", err)
		os.Exit(1)
	}
	client, err := mauth_app.CreateClient(target_url.Scheme + "://" + target_url.Host)
	var response *http.Response
	switch *action {
	case "GET":
		response, err = client.Get(target_url.String())

	case "DELETE":
		response, err = client.Delete(target_url.String())

	case "POST":
		response, err = client.Post(target_url.String(), *data)

	case "PUT":
		response, err = client.Put(target_url.String(), *data)
	}
	defer response.Body.Close()
	fmt.Printf("Status Code: %d\n", response.StatusCode)
	if *headers {
		fmt.Println("Headers:")
		for key, value := range response.Header {
			fmt.Printf(" %s: %s\n", key, value)
		}
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Printf("Response Body:\n%s\n", body)
}
