package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mdsol/go-mauth-client/go_mauth_client"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
)

/*
This uses the underlying library to build a MAuth Client tool;
To build and install:
$ go install
*/

// Context for the MAuth Client
type ApplicationContext struct {
	app_uuid         string
	private_key_text []byte
	private_key_file string
}

// Check that the passed verb is one we prepared for
func CheckAction(action *string) bool {
	switch *action {
	case "GET", "POST", "PUT", "DELETE":
		return true
	default:
		return false
	}
}

func PrettyJson(in []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in, err
	}
	return out.Bytes(), nil
}

// Process the Configuration Content
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

	// No key, textual or file-based, passed in
	if IsNull(&private_key_file) && IsNull(&private_key_text) {
		return nil, errors.New("Need a key specified")
	}
	// Load from text
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

// LoadMAuthConfig loads the Configuration from a JSON file (No YAML support in corelib)
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

// IsNull is a Convenience Function for identification of empty strings
func IsNull(value *string) bool {
	return *value == ""
}

// Main function, go!
func main() {
	// config_file is the path to the configuration file
	config_file := flag.String("config", "", "Specify the configuration file")
	// key_file is the path to the private key file
	key_file := flag.String("private-key", "", "Specify the private key file")
	// app_uuid is the assigned MAuth App ID
	app_uuid := flag.String("app-uuid", "", "Specify the App UUID")
	// action is the HTTP Verb to use on the URL
	action := flag.String("method", "GET", "Specify the method (GET, POST, PUT, DELETE)")
	// data is the data to be POST or PUT
	data := flag.String("data", "", "Specify the data")
	// headers is a flag, which tells the app to print out the response headers
	headers := flag.Bool("headers", false, "Print the Response Headers")
	// verbose is a flag, which tells the app to print out more information
	verbose := flag.Bool("verbose", false, "Print out more information")
	// pretty is a flag, which tells the app to format json
	pretty := flag.Bool("pretty", false, "Prettify the JSON")

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
			log.Fatal("Error loading configuration: ", err)
			os.Exit(1)
		}
	} else {
		var err error
		mauth_app, err = go_mauth_client.LoadMauth(*app_uuid, *key_file)
		if err != nil {
			log.Fatal("Error loading configuration: ", err)
			os.Exit(1)
		}
	}
	if *verbose {
		log.Println("Created MAuth App with App UUID: ", mauth_app.App_ID)
	}
	action_matches := CheckAction(action)
	if !action_matches {
		log.Fatal("Action ", action, "is not known")
		flag.Usage()
		os.Exit(1)
	}
	var args []string
	args = flag.Args()
	if len(args) == 0 {
		log.Fatal("No URL, nothing to do")
		flag.Usage()
		os.Exit(1)
	}
	target_url, err := url.Parse(args[0])
	if err != nil {
		log.Fatal("Unable to parse url: ", err)
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
	if err != nil {
		log.Fatal("Error raised in request ", err, " please check")
	}
	defer response.Body.Close()
	if *verbose {
		log.Println("Status Code: ", response.StatusCode)
	}
	if *headers {
		log.Println("Headers:")
		for key, value := range response.Header {
			log.Printf(" %s: %s\n", key, value)
		}
	}
	body, err := ioutil.ReadAll(response.Body)
	if *verbose {
		log.Println("Response Body:")
	}
	if *pretty {
		media_type, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
		if err == nil {
			if media_type == "application/json" {
				pretty, err := PrettyJson(body)
				if err != nil {
					fmt.Println(string(body))
				} else {
					fmt.Print(string(pretty))
				}
			} else {
				fmt.Println(string(body))
			}
		}

	} else {
		fmt.Println(string(body))
	}
}
