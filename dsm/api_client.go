// **********
// Terraform Provider - DSM: api client
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.5
//       - Date:      27/11/2020
// **********

package dsm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
func (obj *api_client) APICallBody(method string, url string, body map[string]interface{}) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 600 * time.Second}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), bytes.NewBuffer(reqBody))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
		})
	} else {
		req.Header.Add("Authorization", "Bearer "+obj.authtoken)

		r, err := client.Do(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK]: Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
			})
		} else {
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
				})
			} else {
				resp := make(map[string]interface{})
				if r.StatusCode > 204 || r.StatusCode < 200 {
					err = json.Unmarshal(bodyBytes, &resp)
					if err != nil {
						bodyString := string(bodyBytes)
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "[DSM SDK]: Call DSM provider API returned non-JSON",
							Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, bodyString),
						})
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "[DSM SDK]: Call DSM provider API returned error",
							Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, resp),
						})
					}
				} else {
					err = json.Unmarshal(bodyBytes, &resp)
					if err != nil {
						bodyString := string(bodyBytes)
						resp = map[string]interface{}{
							"msg": bodyString,
						}
					}
					return resp, nil
				}
			}
		}
	}
	return nil, diags
}

// [-]: call api without body
func (obj *api_client) APICall(method string, url string) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 600 * time.Second}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
		})
	} else {
		req.Header.Add("Authorization", "Bearer "+obj.authtoken)

		r, err := client.Do(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK]: Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
			})
		} else {
			defer r.Body.Close()

			// FIXME: DELETE does not have any output
			if method == "DELETE" {
				return nil, nil
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
				})
			} else {
				resp := make(map[string]interface{})
				if r.StatusCode > 204 || r.StatusCode < 200 {
					err = json.Unmarshal(bodyBytes, &resp)
					if err != nil {
						bodyString := string(bodyBytes)
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "[DSM SDK]: Call DSM provider API returned non-JSON",
							Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, bodyString),
						})
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "[DSM SDK]: Call DSM provider API returned error",
							Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, resp),
						})
					}
				} else {
					err = json.Unmarshal(bodyBytes, &resp)
					if err != nil {
						bodyString := string(bodyBytes)
						resp = map[string]interface{}{
							"msg": bodyString,
						}
					}
					return resp, nil
				}
			}
		}
	}
	return nil, diags
}

// [-]: call api without body - return as array
func (obj *api_client) APICallList(method string, url string) ([]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 600 * time.Second}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
		})
	} else {
		req.Header.Add("Authorization", "Bearer "+obj.authtoken)

		r, err := client.Do(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK]: Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
			})
		} else {

			defer r.Body.Close()

			// FIXME: DELETE does not have any output
			if method == "DELETE" {
				return nil, nil
			}

			resp := make([]interface{}, 0)
			err = json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
				})
			} else {
				return resp, nil
			}
		}
	}
	return nil, diags
}
