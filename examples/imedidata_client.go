// Package examples is used for examples of using the mauth client library
package examples

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mdsol/go-mauth-client"
)

/*
This is an example for querying data from the iMedidata API using the go-mauth-client
*/

// User struct returned by iMedidata
type User struct {
	Login                    string
	Email                    string
	Uuid                     string
	Activation_code          string
	Activation_status        string
	Name                     string
	First_name               string
	Last_name                string
	Time_zone                string
	Locale                   string
	Institution              string
	Title                    string
	Department               string
	Address_line_1           string
	Address_line_2           string
	Address_line_3           string
	City                     string
	State                    string
	Postal_code              string
	Country                  string
	Telephone                string
	Fax                      string
	Pager                    string
	Mobile                   string
	Creator_uuid             string
	Admin                    string
	Clinical_data_restricted string
}

type UserResponse struct {
	User User
}

// Example implementing:
// http://developer.imedidata.com/desktop/ActionTopics/Users/Listing_User_Account_Details.htm
func getUserDetails(mauthApp *go_mauth_client.MAuthApp, userUuid string) (user *User, err error) {
	client, err := mauthApp.CreateClient("https://innovate.imedidata.com")
	if err != nil {
		log.Fatal("Error creating client")
		return nil, err
	}
	userDetailsResponse, err := client.Get("api/v2/users/" + userUuid + ".json")
	if err != nil {
		log.Fatal("Error downloading User Details")
		return nil, err
	}
	if userDetailsResponse.StatusCode != 200 {
		log.Fatal("Request status code: ", userDetailsResponse.StatusCode)
		return nil, err
	}
	defer userDetailsResponse.Body.Close()

	content, err := ioutil.ReadAll(userDetailsResponse.Body)
	if err != nil {
		log.Fatal("Unable to read response")
		return nil, err
	}
	var userResponse UserResponse

	err = json.Unmarshal(content, &userResponse)
	if err != nil {
		log.Fatal("Unable to deserialise response")
		return nil, err
	}
	user = &userResponse.User
	return user, nil
}

func main() {
	userUUID := os.Getenv("USER_UUID")
	mauthApp, err := loadApp()
	if err != nil {
		log.Fatal("Error creating the client")
	}
	user, err := getUserDetails(mauthApp, userUUID)
	if err != nil {
		log.Fatal(fmt.Printf("Unable to get User: %v", err))
	}
	fmt.Println("User: ", user)
}
