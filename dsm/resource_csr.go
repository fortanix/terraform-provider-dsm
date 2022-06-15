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
	"math/rand"

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
			"st": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"san": {
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

	newsigner, err := NewDSMSigner(d.Get("kid").(string), d.Get("san").(string), d.Get("email").(string), d.Get("cn").(string), d.Get("ou").(string), d.Get("l").(string), d.Get("c").(string), d.Get("o").(string), d.Get("st").(string), m.(*api_client))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to create signer",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %s", newsigner),
		})
		return diags
	}

	generated_csr, err := newsigner.generate_csr()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to get DSM signer up",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %s", err),
		})
		return diags
	}

	if err := d.Set("value", generated_csr); err != nil {
		return diag.FromErr(err)
	}

	idSet := rand.Intn(99999999)
	d.SetId(fmt.Sprintf("%s", idSet))
	return nil
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
