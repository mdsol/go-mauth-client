// Package examples is used for examples of using the mauth client library
package examples

import (
	"github.com/mdsol/go-mauth-client/go_mauth_client"
	"log"
	"encoding/json"
	"io/ioutil"
)

/*
This is an example for querying data from the iMedidata API using the go-mauth-client
 */

// User struct returned by iMedidata
type User struct {
	Login string
	Email string
	Uuid string
	Activation_code string
	Activation_status string
	Name string
	First_name string
	Last_name string
	Time_zone string
	Locale string
	Institution string
	Title string
	Department string
	Address_line_1 string
	Address_line_2 string
	Address_line_3 string
	City string
	State string
	Postal_code string
	Country string
	Telephone string
	Fax string
	Pager string
	Mobile string
	Creator_uuid string
	Admin string
	Clinical_data_restricted string
}

type UserResponse struct {
	User User
}

// Example implementing:
// http://developer.imedidata.com/desktop/ActionTopics/Users/Listing_User_Account_Details.htm
func GetUserDetails(mauth_app *go_mauth_client.MAuthApp, user_uuid string)(user User, err error){
	client, err := mauth_app.CreateClient("https://innovate.imedidata.com")
	if err != nil {
		log.Fatal("Error creating client")
		return nil, err
	}
	user_details_response, err := client.Get("api/v2/users/"+user_uuid+".json")
	if err != nil {
		log.Fatal("Error downloading User Details")
		return nil, err
	}
	if user_details_response.StatusCode != 200 {
		log.Fatal("Request status code: ", user_details_response.StatusCode)
		return nil, err
	}
	defer user_details_response.Body.Close()

	content, err := ioutil.ReadAll(user_details_response.Body)
	if err != nil {
		log.Fatal("Unable to read response")
		return nil, err
	}
	var user_response UserResponse

	err = json.Unmarshal(content, &user_response)
	if err != nil {
		log.Fatal("Unable to deserialise response")
		return nil, err
	}
	user = user_response.User
	return
}
