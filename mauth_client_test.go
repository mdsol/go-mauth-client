package go_mauth_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestFullURLWithRelative(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	expected := "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
	// now, with a trailing slash
	client, _ = mauthApp.CreateClient("https://innovate.mdsol.com/")
	expected = "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ = client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen: ", actual)

	}
}

func TestFullURLWithRelativeAndParams(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	expected := "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
	// now, with a trailing slash
	client, _ = mauthApp.CreateClient("https://innovate.mdsol.com/")
	expected = "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ = client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen: ", actual)

	}
}

func TestFullURLWithActualURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	expected := "https://balance-innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("https://balance-innovate.mdsol.com/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
}

func TestCreateClient(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	if client.baseURL.String() != "https://innovate.mdsol.com" {
		t.Error("Base URL has changed")
	}
	if client.mauthApp.AppID != appID {
		t.Error("App ID has changed")
	}
}

func TestCreateClientBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	_, err := mauthApp.CreateClient("some_nonsense")
	if err == nil {
		t.Error("Bad URL should fail")
	}
}

func TestMauthClient_fullURLBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	_, err := client.fullURL("http://\x7finnovate.mdsol.com")
	if err == nil {
		t.Error("Expected error with Bad URL")
	}
}
func TestMauthClient_fullURLPathBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	_, err := mauthApp.CreateClient("https://\x7finnovate.mdsol.com")
	if err == nil {
		t.Error("Expected error with Bad URL")
	}
}

func TestMauthClient_fullURLPath(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	fullURL, err := client.fullURL("/subject/city")
	if err != nil {
		t.Error("Unexpected error with Path in URL")
	}
	if "https://innovate.mdsol.com/subject/city" != fullURL {
		t.Error("Expected URL generated for path not found")
	}
}

func hasMWSHeader(r *http.Request) bool {
	for header := range r.Header {
		if header == "X-Mws-Authentication" {
			return true
		}
	}
	return false
}

// Test the Get call
func TestMAuthClient_Get(t *testing.T) {
	var verb string
	hasMwsHeader := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.String()
		verb = r.Method
		hasMwsHeader = hasMWSHeader(r)
		w.Header().Set("Content-Type", "application/json")
		// Don't care about errors here
		_, _ = fmt.Fprintln(w, `{"fake twitter json string"}`)

	}))
	defer server.Close()
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient(server.URL)
	// Make the Get call
	_, err := client.Get("/api/v2/users.json")
	if err != nil {
		t.Error("Get Failed: ", err)
	}
	if verb != "GET" {
		t.Error("Expected GET, got ", verb)
	}
	if !hasMwsHeader {
		t.Error("Expected header not present")
	}
}

// Get with bad URL
func TestMAuthClient_GetBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	// Make the Get call
	_, err := client.Get("https://innovate.mdsol.com/api/v2/\x7fusers.json")
	if err == nil {
		t.Error("Expected error with GET to bad URL")
	}
}

// Post with bad URL
func TestMAuthClient_PostBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	// Make the Post call
	_, err := client.Post("https://innovate.mdsol.com/api/v2/\x7fusers.json", "")
	if err == nil {
		t.Error("Expected error with POST to bad URL")
	}
}

// Put with bad URL
func TestMAuthClient_PutBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	// Make the Post call
	_, err := client.Put("https://innovate.mdsol.com/api/v2/\x7fusers.json", "")
	if err == nil {
		t.Error("Expected error with PUT to bad URL")
	}
}

// Put with bad URL
func TestMAuthClient_DeleteBadURL(t *testing.T) {
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient("https://innovate.mdsol.com")
	// Make the Post call
	_, err := client.Delete("https://innovate.mdsol.com/api/v2/\x7fusers.json")
	if err == nil {
		t.Error("Expected error with DELETE to bad URL")
	}
}

// Test the Delete call
func TestMAuthClient_Delete(t *testing.T) {
	var verb string
	hasMwsHeader := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.String()
		verb = r.Method

		hasMwsHeader = hasMWSHeader(r)
		w.Header().Set("Content-Type", "application/json")
		// don't care about errors here
		_, _ = fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))
	defer server.Close()
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient(server.URL)
	// Make the Get call
	_, err := client.Delete("/api/v2/users.json")
	if err != nil {
		t.Error("Delete Failed: ", err)
	}
	if verb != "DELETE" {
		t.Error("Expected DELETE, got ", verb)
	}
	if !hasMwsHeader {
		t.Error("Expected header not present")
	}
}

// Test the Post call
func TestMAuthClient_Post(t *testing.T) {
	var verb string
	hasMwsHeader := false
	var contentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.String()
		verb = r.Method

		hasMwsHeader = hasMWSHeader(r)
		for header, value := range r.Header {
			if header == "Content-Type" {
				contentType = strings.Join(value, "")
			}
		}
		w.Header().Set("Content-Type", "application/json")
		// don't care about error here
		_, _ = fmt.Fprint(w, "{\"fake twitter json string\"}")
	}))
	defer server.Close()
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient(server.URL)
	// Make the Get call
	response, err := client.Post("/api/v2/users.json", `{"uuid":"1234-1234"}`)
	if err != nil {
		t.Error("Post Failed: ", err)
	}
	if verb != "POST" {
		t.Error("Expected POST, got ", verb)
	}
	if !hasMwsHeader {
		t.Error("Expected header not present")
	}
	if contentType != "application/json" {
		t.Error("Expected Content-type not set")
	}
	content, _ := ioutil.ReadAll(response.Body)
	// don't care about error here
	_ = response.Body.Close()
	if strings.Compare(string(content[:]),
		"{\"fake twitter json string\"}") != 0 {
		t.Error("Unexpected response body: ", string(content[:]))
	}
}

// Test the Put call
func TestMAuthClient_Put(t *testing.T) {
	var verb string
	hasMwsHeader := false
	var contentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.String()
		verb = r.Method

		hasMwsHeader = hasMWSHeader(r)
		for header, value := range r.Header {
			if header == "Content-Type" {
				contentType = strings.Join(value, "")
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{\"fake twitter json string\"}")
	}))
	defer server.Close()
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient(server.URL)
	// Make the Get call
	response, err := client.Put("/api/v2/users.json", `{"uuid":"1234-1234"}`)
	if err != nil {
		t.Error("Post Failed: ", err)
	}
	if verb != "PUT" {
		t.Error("Expected PUT, got ", verb)
	}
	if !hasMwsHeader {
		t.Error("Expected header not present")
	}
	if contentType != "application/json" {
		t.Error("Expected Content-type not set")
	}
	content, _ := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if strings.Compare(string(content[:]),
		"{\"fake twitter json string\"}") != 0 {
		t.Error("Unexpected response body: ", string(content[:]))
	}
}

// Test adding a Header
func TestMAuthClient_SetHeader(t *testing.T) {
	var verb string
	hasMwsHeader := false
	hasMccVersionHeader := false
	var contentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.String()
		verb = r.Method

		hasMwsHeader = hasMWSHeader(r)
		for header, value := range r.Header {
			if header == "Content-Type" {
				contentType = strings.Join(value, "")
			}
			if header == "Mcc-Version" {
				if value[0] == "v2019-03-22" {
					hasMccVersionHeader = true
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{\"fake twitter json string\"}")
	}))
	defer server.Close()
	mauthApp, _ := LoadMauth(appID, filepath.Join("test", "private_key.pem"))
	client, _ := mauthApp.CreateClient(server.URL)
	// Set a Header
	client.SetHeader("Mcc-Version", "v2019-03-22")
	// Make the Get call
	response, err := client.Put("/api/v2/users.json", `{"uuid":"1234-1234"}`)
	if err != nil {
		t.Error("Post Failed: ", err)
	}
	if verb != "PUT" {
		t.Error("Expected PUT, got ", verb)
	}
	if !hasMwsHeader {
		t.Error("Expected MWS header not present")
	}
	if !hasMccVersionHeader {
		t.Error("Expected Mcc-version header not present")
	}
	if contentType != "application/json" {
		t.Error("Expected Content-type not set")
	}
	content, _ := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if strings.Compare(string(content[:]),
		"{\"fake twitter json string\"}") != 0 {
		t.Error("Unexpected response body: ", string(content[:]))
	}

}

// Example of creating a MAuth Client
func ExampleMAuthApp_CreateClient() {
	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth mAuthApp
	var mAuthApp *MAuthApp
	mAuthApp, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create mAuthApp: ", err)
	}
	// Define a base URL
	var baseURL = "https://innovate.imedidata.com"
	var client *MAuthClient
	client, err = mAuthApp.CreateClient(baseURL)
	if err != nil {
		log.Fatal("Unable to create MAuth Client: ", err)
	}
	println("Successfully created MAuth Client for APP: ", client.mauthApp.AppID)
}

// Example of creating a MAuth Client and making a Get Request
func ExampleMAuthClient_Get() {
	// Get information on a User
	// http://developer.imedidata.com/desktop/ActionTopics/Users/Listing_User_Account_Details.htm

	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth mAuthApp
	var mAuthApp *MAuthApp
	mAuthApp, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create mAuthApp: ", err)
	}
	// Define a base URL
	var baseURL = "https://innovate.imedidata.com"

	// Define and create the Client
	var client *MAuthClient
	client, err = mAuthApp.CreateClient(baseURL)
	if err != nil {
		log.Fatal("Unable to create MAuth Client: ", err)
	}
	// This is made-up
	var userUuid = "347942BF-9915-405D-BB20-6196597F3BE3"
	response, err := client.Get("api/v2/users/" + userUuid + ".json")
	println("Got a status code of", response.StatusCode, "for request for User UUID", userUuid)
}

func ExampleMAuthClient_Post() {
	// Creating a Study Using a MAuth Client
	// http://developer.imedidata.com/desktop/ActionTopics/Studies/Creating_Studies.htm

	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth mAuthApp
	var mAuthApp *MAuthApp
	mAuthApp, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create mAuthApp: ", err)
	}
	// Define a base URL
	var baseURL = "https://innovate.imedidata.com"

	// Define and create the Client
	client, err := mAuthApp.CreateClient(baseURL)
	if err != nil {
		log.Fatal("Unable to create MAuth Client: ", err)
	}

	// Define the constituent entity references
	var studyGroupUUID = "347942BF-9915-405D-BB20-6196597F3BE3"
	var studyUUID = "C3C79E4A-4BFD-4A72-89E9-724A4E6A9D95"

	// This is a slimmed down version of the structure from the reference above
	type studyDefinition struct {
		Number           int    `json:"number"`
		Name             string `json:"name"`
		IsProduction     bool   `json:"is_production"`
		TherapeticArea   string `json:"therapeutic_area"`
		FullDescription  string `json:"full_description"`
		CompoundCode     string `json:"compound_code"`
		DrugDevice       string `json:"drug_device"`
		Title            string `json:"title"`
		UUID             string
		Protocol         string `json:"protocol"`
		ParentUUID       string `json:"parent_UUID"`
		EnrollmentTarget int    `json:"enrollment_target"`
		OID              string `json:"oid"`
	}

	// Create an instance of the new study
	study := &studyDefinition{
		Number:           1,
		Name:             "ABC1234",
		IsProduction:     true,
		TherapeticArea:   "Endocrine",
		FullDescription:  "Some Sample Study",
		CompoundCode:     "Mediflex",
		DrugDevice:       "Drug",
		Title:            "A sample Endocrine Study",
		UUID:             studyUUID,
		Protocol:         "ABC1234",
		ParentUUID:       "",
		EnrollmentTarget: 150,
		OID:              "ABC1234",
	}
	data, _ := json.Marshal(study)

	// POST www.imedidata.com/api/v2/study_groups/[study group uuid]/studies.json
	response, err := client.Post("api/v2/study_groups/"+studyGroupUUID+"/studies.json",
		string(data))
	println("Got a status code of", response.StatusCode, "for request to create Study", studyUUID)
}
