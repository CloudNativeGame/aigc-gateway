package user

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type CustomData map[string]interface{}

type PatchBody struct {
	CustomData CustomData `json:"customData"`
}

type tokenAuth struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type User struct {
	Id            string                 `json:"id"`
	Username      string                 `json:"username"`
	PrimaryEmail  string                 `json:"primaryEmail"`
	PrimaryPhone  string                 `json:"primaryPhone"`
	Name          string                 `json:"name"`
	Avatar        string                 `json:"avatar"`
	CustomData    map[string]interface{} `json:"customData"`
	LastSignInAt  int                    `json:"lastSignInAt"`
	CreatedAt     int                    `json:"createdAt"`
	ApplicationId string                 `json:"applicationId"`
	IsSuspended   bool                   `json:"isSuspended"`
	HasPassword   bool                   `json:"hasPassword"`
}

var M2MID = os.Getenv("M2M_Id")
var M2MSecret = os.Getenv("M2M_Secret")
var Endpoint = os.Getenv("Endpoint")

func UpdateUserMetaData(userId string, customData map[string]interface{}) error {

	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("resource", "https://default.logto.app/api")
	form.Add("scope", "all")

	tokenRequest, err := http.NewRequest(http.MethodPost, Endpoint+"oidc/token", strings.NewReader(form.Encode()))

	if err != nil {
		return err
	}
	tokenRequest.SetBasicAuth(M2MID, M2MSecret)
	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(tokenRequest)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	au := &tokenAuth{}
	err = json.Unmarshal(b, au)

	if err != nil {
		return err
	}

	patchBody := &PatchBody{}
	patchBody.CustomData = customData

	customDataBytes, err := json.Marshal(patchBody)

	patchRequest, err := http.NewRequest(http.MethodPatch, fmt.Sprintf(Endpoint+"api/users/%s/custom-data", userId), bytes.NewBuffer(customDataBytes))
	if err != nil {
		return err
	}
	patchRequest.Header.Set("Authorization", "Bearer "+au.AccessToken)
	patchRequest.Header.Set("Content-Type", "application/json")

	patchResponse, err := c.Do(patchRequest)
	defer patchResponse.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func GetUserProfile(userId string) (*User, error) {
	c := http.Client{}
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("resource", "https://default.logto.app/api")
	form.Add("scope", "all")

	tokenRequest, err := http.NewRequest(http.MethodPost, Endpoint+"oidc/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	tokenRequest.SetBasicAuth(M2MID, M2MSecret)
	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(tokenRequest)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	au := &tokenAuth{}
	err = json.Unmarshal(b, au)

	if err != nil {
		return nil, err
	}

	userResp, err := http.Get(Endpoint + "api/users/" + userId)

	defer userResp.Body.Close()

	if err != nil {
		return nil, err
	}
	userBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	u := &User{}

	err = json.Unmarshal(userBytes, u)

	return u, err

}
