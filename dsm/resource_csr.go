// **********
// Terraform Provider - DSM: resource: csr
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.15
//       - Date:      26/05/2022
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Security Object
func resourceCsr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateCsr,
		ReadContext:   resourceReadCsr,
		UpdateContext: resourceUpdateCsr,
		DeleteContext: resourceDeleteCsr,
		Schema: map[string]*schema.Schema{
			"kid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ou": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"o": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"l": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"c": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"cn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create CSR
func resourceCreateCsr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subject_dn := make(map[string]interface{})

	if cn := d.Get("cn").(string); len(cn) > 0 {
		subject_dn["CN"] = cn
	}

	if ou := d.Get("ou").(string); len(ou) > 0 {
		subject_dn["OU"] = ou
	}

	if l := d.Get("l").(string); len(l) > 0 {
		subject_dn["L"] = l
	}

	if c := d.Get("c").(string); len(c) > 0 {
		subject_dn["C"] = c
	}

	if o := d.Get("o").(string); len(o) > 0 {
		subject_dn["O"] = o
	}

	plugin_object := map[string]interface{}{
		"subject_key": d.Get("kid").(string),
		"subject_dn":  subject_dn,
	}

	reqfpi, err := m.(*api_client).FindPluginId("Terraform Plugin - CSR")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/plugins: %v", err),
		})
		return diags
	}
	var endpoint = fmt.Sprintf("sys/v1/plugins/%s", string(reqfpi))
	var operation = "POST"

	req, err := m.(*api_client).APICallBody(operation, endpoint, plugin_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/plugins: %v", err),
		})
		return diags
	}

	if err := d.Set("kid", req["kid"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("id", req["id"].(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", req["value"].(string)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("id").(string))
	return resourceReadCsr(ctx, d, m)
}

// [R]: Read Security Object
func resourceReadCsr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [U]: Update Security Object
func resourceUpdateCsr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Security Object
func resourceDeleteCsr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
