// **********
// Terraform Provider - DSM: api client
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.7
//       - Date:      27/11/2020
// **********

package dsm

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"golang.org/x/time/rate"
)

type api_client struct {
	endpoint     string
	port         int
	authtoken    string
	authtype     string
	acct_id      string
	aws_profile  string
	aws_region   string
	azure_region string
	insecure     bool
	timeout      int
}

type dsm_plugin struct {
	Id   string `json:"plugin_id"`
	Name string `json:"name"`
}

// FXHTTPClient Rate Limited HTTP Client
type FXHTTPClient struct {
	client      *retryablehttp.Client
	Ratelimiter *rate.Limiter
}

// [-]: do FXHTTPClient
func (c *FXHTTPClient) Do(req *retryablehttp.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// [-]: set FXHTTPClient
func NewClient(rl *rate.Limiter) *FXHTTPClient {
	c := &FXHTTPClient{
		client:      retryablehttp.NewClient(),
		Ratelimiter: rl,
	}
	return c
}

// [-]: set api_client state
func NewAPIClient(endpoint string, port int, username string, password string, api_key string, acct_id string, aws_profile string, aws_region string, azure_region string, insecure bool, timeout int) (*api_client, error) {
	// FIXME: clunky way of creating api_client session
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	rl := rate.NewLimiter(rate.Every(1*time.Second), 5) // 5 requests in every second
	client := NewClient(rl)
	client.client.HTTPClient.Transport = tr
	client.client.HTTPClient.Timeout = time.Duration(timeout) * time.Second

	resp := make(map[string]interface{})
	var authtype string
	var authtoken string

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/sys/v1/session/auth", endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	if len(api_key) > 0 {
		authtype = "Basic "
		authtoken = api_key
		req.Header.Add("Authorization", authtype+authtoken)
	} else if len(username) > 0 && len(password) > 0 {
		req.SetBasicAuth(username, password)
	} else {
		return nil, fmt.Errorf("unauthorized access to DSM")
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 401 {
		return nil, fmt.Errorf("unauthorized access to DSM")
	}

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

	req, err = retryablehttp.NewRequest("POST", fmt.Sprintf("%s/sys/v1/session/select_account", endpoint), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	authtype = "Bearer "
	authtoken = resp["access_token"].(string)
	req.Header.Add("Authorization", authtype+authtoken)
	req.Close = true

	r, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// EOF error: select_acccount has no return
	_, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

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

				req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/sys/v1/session/aws_temporary_credentials", endpoint), bytes.NewBuffer(reqBody))
				if err != nil {
					return nil, err
				}
				req.Header.Add("Authorization", "Bearer "+resp["access_token"].(string))
				req.Close = true

				r, err := client.Do(req)
				if err != nil {
					return nil, err
				}
				defer r.Body.Close()

				// EOF error: aws_temporary_credentials has no return
				_, err = io.ReadAll(r.Body)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	newclient := api_client{
		endpoint:     endpoint,
		port:         port,
		authtoken:    authtoken,
		authtype:     authtype,
		acct_id:      acct_id,
		aws_profile:  aws_profile,
		aws_region:   aws_region,
		azure_region: azure_region,
		insecure:     insecure,
		timeout:      timeout,
	}
	return &newclient, nil
}

// [-]: call api with body
func (obj *api_client) APICallBody(method string, url string, body map[string]interface{}) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	rl := rate.NewLimiter(rate.Every(1*time.Second), 5) // 5 requests in every second
	client := NewClient(rl)
	client.client.HTTPClient.Transport = tr
	client.client.HTTPClient.Timeout = time.Duration(obj.timeout) * time.Second

	reqBody, err := json.MarshalIndent(&body, "", "\t")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[Conversion] Unable to marshal request body",
			Detail:   fmt.Sprintf("[Conversion] Body: %s - Err: %v", body, err),
		})
		return nil, diags
	}

	req, err := retryablehttp.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), bytes.NewBuffer(reqBody))
	req.Close = true
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
		})
	} else {
		req.Header.Add("Authorization", obj.authtype+obj.authtoken)

		r, err := client.Do(req)
		if err != nil {
			if r != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to call DSM provider API client",
					Detail:   fmt.Sprintf("[E]: API: %s %s %d: %s", method, url, r.StatusCode, err),
				})
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to call DSM provider API client",
					Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
				})
			}
		} else {
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s %d: %s", method, url, r.StatusCode, err),
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
							Detail:   fmt.Sprintf("[E]: API: %s %s %d: %s", method, url, r.StatusCode, bodyString),
						})
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "[DSM SDK]: Call DSM provider API returned error",
							Detail:   fmt.Sprintf("[E]: API: %s %s %d: %s", method, url, r.StatusCode, resp),
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
func (obj *api_client) APICall(method string, url string) (map[string]interface{}, int, diag.Diagnostics) {
	var diags diag.Diagnostics
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	rl := rate.NewLimiter(rate.Every(1*time.Second), 5) // 5 requests in every second
	client := NewClient(rl)
	client.client.HTTPClient.Transport = tr
	client.client.HTTPClient.Timeout = time.Duration(obj.timeout) * time.Second

	req, err := retryablehttp.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), nil)
	req.Close = true
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
		})
	} else {
		req.Header.Add("Authorization", obj.authtype+obj.authtoken)

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
				// Only check status
				if r.StatusCode == 204 {
					return nil, r.StatusCode, nil
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "[DSM SDK]: DSM provider API call failed",
						Detail:   fmt.Sprintf("[E]: API: %s %s: %d", method, url, r.StatusCode),
					})
					return nil, r.StatusCode, diags
				}
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: %s %s %d: %s", method, url, r.StatusCode, err),
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
							Detail:   fmt.Sprintf("[E]: API: %s %s %d: %s", method, url, r.StatusCode, bodyString),
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
					return resp, r.StatusCode, nil
				}
			}
		}
	}
	return nil, 500, diags
}

// [-]: call api without body - return as array
func (obj *api_client) APICallList(method string, url string) ([]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: obj.insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	rl := rate.NewLimiter(rate.Every(1*time.Second), 5) // 5 requests in every second
	client := NewClient(rl)
	client.client.HTTPClient.Transport = tr
	client.client.HTTPClient.Timeout = time.Duration(obj.timeout) * time.Second

	req, err := retryablehttp.NewRequest(method, fmt.Sprintf("%s/%s", obj.endpoint, url), nil)
	req.Close = true
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: %s %s: %s", method, url, err),
		})
	} else {
		req.Header.Add("Authorization", obj.authtype+obj.authtoken)

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
		Proxy:           http.ProxyFromEnvironment,
	}

	rl := rate.NewLimiter(rate.Every(1*time.Second), 5) // 5 requests in every second
	client := NewClient(rl)
	client.client.HTTPClient.Transport = tr
	client.client.HTTPClient.Timeout = time.Duration(obj.timeout) * time.Second

	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/sys/v1/plugins", obj.endpoint), nil)
	req.Close = true
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to prepare DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %s", err),
		})
	} else {
		req.Header.Add("Authorization", obj.authtype+obj.authtoken)

		r, err := client.Do(req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK]: Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %s", err),
			})
		} else {
			defer r.Body.Close()

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "[DSM SDK]: Unable to read DSM provider API response",
					Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %s", err),
				})
			} else {
				var allPlugins []dsm_plugin
				err = json.Unmarshal(bodyBytes, &allPlugins)
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "[DSM SDK]: Unable to unmarshal provider API response body",
						Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %s", err),
					})
				}
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
						Detail:   "[E]: API: GET sys/v1/plugins",
					})
				} else {
					return []byte(resp), nil
				}
			}
		}
	}
	return nil, diags
}
