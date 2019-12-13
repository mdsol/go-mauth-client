// go_mauth_client is a small application using the MAuth Library to make signed calls against Medidata APIs
//
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/go-xmlfmt/xmlfmt"
	go_mauth_client "github.com/mdsol/go-mauth-client"
)

// Context for the MAuth Client
type ApplicationContext struct {
	appUuid        string
	privateKeyText []byte
	privateKeyFile string
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

// Prettify JSON Output
func PrettyJson(in []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in, err
	}
	return out.Bytes(), nil
}

// Process the Configuration Content
func ProcessConfiguration(content []byte) (mauthApp *go_mauth_client.MAuthApp, err error) {
	var context map[string]string
	err = json.Unmarshal(content, &context)
	if err != nil {
		return nil, err
	}
	appUuid := context["app_uuid"]
	privateKeyFile := context["private_key_file"]
	privateKeyText := context["private_key_text"]
	disableV1, err := strconv.ParseBool(context["disable_v1"])
	if IsNull(&appUuid) {
		return nil, errors.New("Need an app_uuid specified")
	}

	// No key, textual or file-based, passed in
	if IsNull(&privateKeyFile) && IsNull(&privateKeyText) {
		return nil, errors.New("Need a key specified")
	}
	// Load from text
	if IsNull(&privateKeyFile) {
		// read from the embedded value
		mauthApp, err = go_mauth_client.LoadMauth(go_mauth_client.MAuthOptions{appUuid, privateKeyText, disableV1})
	} else {
		// load the key from a file
		mauthApp, err = go_mauth_client.LoadMauth(go_mauth_client.MAuthOptions{appUuid, privateKeyFile, disableV1})
	}
	return
}

// LoadMAuthConfig loads the Configuration from a JSON file (No YAML support in corelib)
func LoadMAuthConfig(fileName string) (mauthApp *go_mauth_client.MAuthApp, err error) {
	if _, err = os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	mauthApp, err = ProcessConfiguration(content)
	return
}

// IsNull is a Convenience Function for identification of empty strings
func IsNull(value *string) bool {
	return *value == ""
}

// Main function, go!
func main() {
	// config_file is the path to the configuration file
	configFile := flag.String("config", "", "Specify the configuration file")
	// key_file is the path to the private key file
	keyFile := flag.String("private-key", "", "Specify the private key file")
	// app_uuid is the assigned MAuth App ID
	appUuid := flag.String("app-uuid", "", "Specify the App UUID")
	// disableV1 specifies if V1 signing should occurr
	disableV1 := flag.Bool("disableV1", false, "Specify if V1 signing should be disabled")
	// action is the HTTP Verb to use on the URL
	action := flag.String("method", "GET", "Specify the method (GET, POST, PUT, DELETE)")
	// data is the data to be POST or PUT
	data := flag.String("data", "", "Specify the data")
	// headers is a flag, which tells the app to print out the response headers
	headers := flag.Bool("headers", false, "Print the Response Headers")
	// verbose is a flag, which tells the app to print out more information
	verbose := flag.Bool("verbose", false, "Print out more information")
	// pretty is a flag, which tells the app to format output (JSON/XML)
	pretty := flag.Bool("pretty", false, "Prettify the Output")
	// Mcc-Version
	mccVersion := flag.String("mcc-version", "", "Specify the MCC version for the endpoint")
	// version is a flag, which tells the app to format json
	version := flag.Bool("version", false, "Print out the version")

	flag.Parse()

	if *version == true {
		println("Go MAuth Client CLI: Version ", go_mauth_client.GetVersion())
		os.Exit(0)
	}
	// No information supplied
	if IsNull(configFile) && (IsNull(keyFile) || IsNull(appUuid)) {
		println("Need to specify configuration file or app settings")
		flag.Usage()
		os.Exit(1)
	}
	var mauthApp *go_mauth_client.MAuthApp
	// Load the MAuth Config
	if !IsNull(configFile) {
		// a config file has been passed
		var err error
		mauthApp, err = LoadMAuthConfig(*configFile)
		if err != nil {
			log.Fatal("Error loading configuration: ", err)
			os.Exit(1)
		}
	} else {
		var err error
		mauthApp, err = go_mauth_client.LoadMauth(go_mauth_client.MAuthOptions{*appUuid, *keyFile, *disableV1})
		if err != nil {
			log.Fatal("Error loading configuration: ", err)
			os.Exit(1)
		}
	}
	if *verbose {
		log.Println("Created MAuth App with App UUID: ", mauthApp.AppId)
	}
	actionMatches := CheckAction(action)
	if !actionMatches {
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
	targetUrl, err := url.Parse(args[0])
	if err != nil {
		log.Fatal("Unable to parse url: ", err)
		os.Exit(1)
	}
	client, err := mauthApp.CreateClient(targetUrl.Scheme + "://" + targetUrl.Host)
	if err != nil {
		log.Fatalf("Error creating MAuthClient: %s", err)
		os.Exit(1)
	}
	if *mccVersion != "" {
		client.SetHeader("Mcc-version", *mccVersion)
	}
	var response *http.Response
	switch *action {
	case "GET":
		response, err = client.Get(targetUrl.String())

	case "DELETE":
		response, err = client.Delete(targetUrl.String())

	case "POST":
		response, err = client.Post(targetUrl.String(), *data)

	case "PUT":
		response, err = client.Put(targetUrl.String(), *data)
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
		mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
		if err == nil {
			if mediaType == "application/json" {
				// Format the JSON output
				pretty, err := PrettyJson(body)
				if err != nil {
					fmt.Println(string(body))
				} else {
					fmt.Print(string(pretty))
				}
			} else if mediaType == "application/xml" || mediaType == "text/xml" {
				// Format the XML output
				fmt.Println(xmlfmt.FormatXML(string(body), "  ", "  "))
			} else {
				fmt.Println(string(body))
			}
		}

	} else {
		fmt.Println(string(body))
	}
}
