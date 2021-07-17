// **********
// Terraform Provider - DSM: resource: aws security object
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.9
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define AWS Security Object
func resourceAWSSobject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAWSSobject,
		ReadContext:   resourceReadAWSSobject,
		UpdateContext: resourceUpdateAWSSobject,
		DeleteContext: resourceDeleteAWSSobject,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acct_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"obj_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key_ops": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [-]: AWS - Load Profile Credentials
func loadAWSProfileCreds(profile_name string, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Specify profile to load for the session's config
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile_name,
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	if sess != nil {
		output, err := sess.Config.Credentials.Get()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[AWS SDK]: Unable to retrieve IAM or STS creds",
				Detail:   fmt.Sprintf("[E]: SDK: AWS credentials access failure: %s", err),
			})
			return diags
		} else {
			aws_temporary_credentials := map[string]interface{}{
				"access_key":    output.AccessKeyID,
				"secret_key":    output.SecretAccessKey,
				"session_token": output.SessionToken,
			}
			_, err := m.(*api_client).APICallBody("POST", "sys/v1/session/aws_temporary_credentials", aws_temporary_credentials)
			if err != nil {
				return err
			}
		}
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[AWS SDK]: Unable to setup session",
			Detail:   fmt.Sprintf("[E]: SDK: AWS session failure: %s", err),
		})
		return diags
	}
	return nil
}

// [C]: Create AWS Security Object
func resourceCreateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Check if AWS Profile is set and use it
	if d.Get("profile") != nil {
		err := loadAWSProfileCreds(d.Get("profile").(string), m)
		if err != nil {
			return err
		}
	}

	security_object := map[string]interface{}{
		"name":        d.Get("name").(string),
		"group_id":    d.Get("group_id").(string),
		"key":         d.Get("key"),
		"description": d.Get("description").(string),
	}

	if err := d.Get("custom_metadata").(map[string]interface{}); len(err) > 0 {
		security_object["custom_metadata"] = d.Get("custom_metadata")
	}

	check_hmg_req := map[string]interface{}{}
	// Scan the AWS Group first before
	req, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/check", d.Get("group_id").(string)), check_hmg_req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/check: %s", err),
		})
		return diags
	}

	req, err = m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/scan", d.Get("group_id").(string)), check_hmg_req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/scan: %s", err),
		})
		return diags
	}

	req, err = m.(*api_client).APICallBody("POST", "crypto/v1/keys/copy", security_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST crypto/v1/keys/copy: %s", err),
		})
		return diags
	}

	d.SetId(req["kid"].(string))
	return resourceReadAWSSobject(ctx, d, m)
}

// [R]: Read AWS Security Object
func resourceReadAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Check if AWS Profile is set and use it
	if d.Get("profile") != nil {
		err := loadAWSProfileCreds(d.Get("profile").(string), m)
		if err != nil {
			return err
		}
	}

	req, err := m.(*api_client).APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET crypto/v1/keys: %s", err),
		})
		return diags
	}

	if err := d.Set("name", req["name"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", req["group_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", req["links"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kid", req["kid"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creator", req["creator"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_metadata", req["custom_metadata"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key_ops", req["key_ops"]); err != nil {
		return diag.FromErr(err)
	}
	if _, ok := req["description"]; ok {
		if err := d.Set("description", req["description"].(string)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("enabled", req["enabled"].(bool)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", req["state"].(string)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// [U]: Update AWS Security Object
func resourceUpdateAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete AWS Security Object
func resourceDeleteAWSSobject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Check if AWS Profile is set and use it
	if d.Get("profile") != nil {
		err := loadAWSProfileCreds(d.Get("profile").(string), m)
		if err != nil {
			return err
		}
	}

	// FIXME: Since deleting, might as well remove the alias
	if d.Get("custom_metadata").(map[string]interface{})["aws-aliases"] != "" {
		remove_aws_alias := map[string]interface{}{
			"kid": d.Id(),
		}
		remove_aws_alias["custom_metadata"] = map[string]interface{}{
			"aws-aliases": "",
			"aws-policy":  d.Get("custom_metadata").(map[string]interface{})["aws-policy"],
		}
		_, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("crypto/v1/keys/%s", d.Id()), remove_aws_alias)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: PATCH crypto/v1/keys: %s", err),
			})
			return diags
		}
	}

	// FIXME: Need to schedule deletion then delete the key - default is set to 7 days for now
	delete_object := map[string]interface{}{
		"pending_window_in_days": 7,
	}
	if d.Get("custom_metadata").(map[string]interface{})["aws-key-state"] != "PendingDeletion" {
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("crypto/v1/keys/%s/schedule_deletion", d.Id()), delete_object)
		if err != nil {
			return err
		}
	}

	_, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("crypto/v1/keys/%s", d.Id()))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE crypto/v1/keys: %s", err),
		})
		return diags
	}

	d.SetId("")
	return nil
}
