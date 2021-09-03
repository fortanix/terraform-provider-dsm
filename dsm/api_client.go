// **********
// Terraform Provider - DSM: api client
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.2.4
//       - Date:      27/11/2020
// **********

package dsm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type api_client struct {
	endpoint    string
	port        int
	authtoken   string
	acct_id     string
	aws_profile string
	aws_region  string
	insecure    bool
}

type dsm_plugin struct {
	Id   string `json:"plugin_id"`
	Name string `json:"name"`
}

// [-]: set api_client state
func NewAPIClient(endpoint string, port int, username string, password string, acct_id string, aws_profile string, aws_region string, insecure bool) (*api_client, error) {
	// FIXME: clunky way of creating api_client session
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	client := &http.Client{Timeout: 600 * time.Second, Transport: tr}

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

	// Check if AWS profile is set and use it within API client
	if len(aws_profile) > 0 {
		// Specify profile to load for the session's config
		sess, err := session.NewSessionWithOptions(session.Options{
			Profile: aws_profile,
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
			// Force enable Shared Config support
			SharedConfigState: session.SharedConfigEnable,
		})
		if sess != nil {
			output, err := sess.Config.Credentials.Get()
			if err != nil {
				return nil, err
			} else {
				aws_temporary_credentials := map[string]interface{}{
					"access_key":    output.AccessKeyID,
					"secret_key":    output.SecretAccessKey,
					"session_token": output.SessionToken,
				}
				reqBody, err := json.Marshal(aws_temporary_credentials)
				if err != nil {
					return nil, err
				}

				req, err = http.NewRequest("POST", fmt.Sprintf("%s/sys/v1/session/aws_temporary_credentials", endpoint), bytes.NewBuffer(reqBody))
				if err != nil {
					return nil, err
				}
				req.Header.Add("Authorization", "Bearer "+resp["access_token"].(string))

				r, err = client.Do(req)
				if err != nil {
					return nil, err
				}
				defer r.Body.Close()
			}
		} else {
			return nil, err
		}
	}

	newclient := api_client{
		endpoint:    endpoint,
		port:        port,
		authtoken:   resp["access_token"].(string),
		acct_id:     acct_id,
		aws_profile: aws_profile,
		aws_region:  aws_region,
		insecure:    insecure,
	}
	return &newclient, nil
}

// [-]: call api with body
func (obj *api_client) APICallBody(method string, url string, body map[string]interface{}) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
	}

	client := &http.Client{Timeout: 600 * time.Second, Transport: tr}
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
				Detail:   fmt.Sprintf("[E]: API: %s %s %s: %s", method, url, r.StatusCode, err),
			})
		} else {
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s %s: %s", method, url, r.StatusCode, err),
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
							Detail:   fmt.Sprintf("[E]: API: %s %s %s: %s", method, url, r.StatusCode, bodyString),
						})
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "[DSM SDK]: Call DSM provider API returned error",
							Detail:   fmt.Sprintf("[E]: API: %s %s %s: %s", method, url, r.StatusCode, resp),
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
	}

	client := &http.Client{Timeout: 600 * time.Second, Transport: tr}

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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
	}

	client := &http.Client{Timeout: 600 * time.Second, Transport: tr}
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

// [-]: find plugin - "Terraform Plugin" - return as array
func (obj *api_client) FindPluginId(plugin_name string) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
	}

	client := &http.Client{Timeout: 60 * time.Second, Transport: tr}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sys/v1/plugins", obj.endpoint), nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", "GET", "sys/v1/plugins", err),
		})
	} else {
		req.Header.Add("Authorization", "Bearer "+obj.authtoken)

		r, err := client.Do(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK]: Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: %s %s: %s", "GET", "sys/v1/plugins", err),
			})
		} else {
			defer r.Body.Close()

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s: %s", "GET", "sys/v1/plugins", err),
				})
			} else {
				var allPlugins []dsm_plugin
				err = json.Unmarshal(bodyBytes, &allPlugins)
				resp := ""
				for i := range allPlugins {
					if allPlugins[i].Name == plugin_name {
						resp = allPlugins[i].Id
						break
					}
				}
				if resp == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "[DSM SDK]: Unable to find Terraform Plugin through DSM provider",
						Detail:   fmt.Sprintf("[E]: API: %s %s", "GET", "sys/v1/plugins"),
					})
				} else {
					return []byte(resp), nil
				}
			}
		}
	}
	return nil, diags
}
