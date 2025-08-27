// **********
// Terraform Provider - SDKMS: resource: security object
// **********
//       - Author:    Ravi Gopal at fortanix dot com
//       - Version:   0.5.36
//       - Date:      24/07/2025
// **********

package dsm

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// [-] Define API
func resourceAPI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAPI,
		ReadContext:   resourceReadAPI,
		UpdateContext: resourceUpdateAPI,
		DeleteContext: resourceDeleteAPI,
		Description: "Triggers the DSM API using the details provided in the request.",
		Schema: map[string]*schema.Schema{
			"method": {
				Description: "HTTP method. Configure one of the parameters below. \n" +
				"   * POST, GET, PATCH, DELETE",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"POST", "PATCH", "GET", "DELETE"}, true),
			},
			"resource_type": {
				Description: "Dsm resource type. Configure one of the parameters below.\n" +
				"   * key, plugin, app, group, user, user_invite",
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_uuid": {
				Description: "DSM resource UUID.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"payload": {
				Description: "JSON body to invoke the request.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_id_attribute": {
				Description: "Name of the response field to be used as the resource's unique identifier.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_response": {
				Description: "Response of the given request.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"recall": {
				Description: "Trigger the API for the same request.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// Create
func resourceCreateAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	req, err := invokeAPI(ctx, d, m)
	if err != nil {
		return err
	}
	// Set the ID for resource
	if api_id_attribute := d.Get("api_id_attribute").(string); len(api_id_attribute) > 0 {
		if api_id, ok := req[api_id_attribute]; ok {
			d.SetId(api_id.(string));
			return nil
		}
	}
	d.SetId(generateRandomID())
	return nil
}

// Read
func resourceReadAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// This is not required.
	return nil
}

// Update
func resourceUpdateAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("method", "resource_type", "payload", "resource_uuid", "recall") {
		_, err := invokeAPI(ctx, d, m)
		return err
	}
	return nil
}

// Delete
func resourceDeleteAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

// Inovke API, it is same for both create and update.
func invokeAPI(ctx context.Context, d *schema.ResourceData, m interface{}) (map[string]interface{}, diag.Diagnostics) {
	http_method := d.Get("method").(string)
	resource_type := d.Get("resource_type").(string)
	// recall is set to false, when a user want recall the same API again,
	// this shows the changes and inovke the API again
	d.Set("recall", false)
	var endpoint string
	endpoint, ok := dsm_endpoints[resource_type]
	if !ok {
		return nil, invokeErrorDiagsWithSummary("Provide the correct value for resource_type", error_summary)
	}
	if resource_uuid := d.Get("resource_uuid").(string); len(resource_uuid) > 0 {
		endpoint += "/" + resource_uuid
	}
	var err diag.Diagnostics
	var req map[string]interface{}
	if http_method == "GET" || http_method == "DELETE" {
		req, _, err = m.(*api_client).APICall(http_method, endpoint)
	} else {
		var payload = map[string]interface{}{}
		if json_body := d.Get("payload").(string); len(json_body) > 0 {
			payload, _ = ConvertStringToJSONGeneric[map[string]interface{}](json_body)
		}
		req, err = m.(*api_client).APICallBody(http_method, endpoint, payload)
	}
	if err != nil {
		return nil, invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: %s %s: %v", http_method, endpoint, err), error_summary)
	}
	reqStr, reqErr := json.Marshal(req)
	if reqErr == nil {
		if err := d.Set("api_response", string(reqStr)); err != nil {
			return nil, diag.FromErr(err)
		}
	} else {
		return nil, invokeErrorDiagsWithSummary(fmt.Sprintf("[E]: API: %s %s: %v", http_method, endpoint, reqErr), error_summary)
	}
	return req, nil
}
