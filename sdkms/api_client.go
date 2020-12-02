// **********
// Terraform Provider - SDKMS: api client
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1
//       - Date:      27/11/2020
// **********

package sdkms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type api_client struct {
	endpoint  string
	port      int
	authtoken string
	acct_id   string
}

// [-]: set api_client state
func NewAPIClient(endpoint string, port int, username string, password string, acct_id string) (*api_client, error) {
	// FIXME: clunky way of creating api_client session
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sys/v1/session/auth", endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	resp := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	account_object := map[string]interface{}{
		"acct_id": acct_id,
	}

	reqBody, err := json.Marshal(account_object)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/sys/v1/session/select_account", endpoint), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+resp["access_token"].(string))

	r, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	newclient := api_client{
		endpoint:  endpoint,
		port:      port,
		authtoken: resp["access_token"].(string),
		acct_id:   acct_id,
	}
	return &newclient, nil
}

// [-]: call api with body
func (obj *api_client) APICallBody(method string, url string, body map[string]interface{}) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+obj.authtoken)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	resp := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// [-]: call api without body
func (obj *api_client) APICall(method string, url string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+obj.authtoken)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// FIXME: DELETE does not have any output
	if method == "DELETE" {
		return nil, nil
	}

	resp := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
