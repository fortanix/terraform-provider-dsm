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
	"net"
	"strconv"

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
		Description: "Generates a CSR from an existing private key within DSM.",
		Schema: map[string]*schema.Schema{
			"kid": {
			    Description: "The security object kid.",
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
			    Description: "The security object value of Generated CSR.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"ou": {
			    Description: "The security object CSR Organisational Unit.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"o": {
			    Description: "The security object CSR Organisation.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"l": {
			    Description: "The security object CSR Location.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"c": {
			    Description: "The security object CSR Country.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"st": {
			    Description: "The security object CSR State.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"e": {
			    Description: "The security object CSR Email.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"email": {
			    Description: "Email value for CSR in Subject Alternative names.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dnsnames": {
			    Description: "The security object CSR DNS Names.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ips": {
			    Description: "The security object CSR IPs.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cn": {
			    Description: "The security object CSR Common Name.",
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
			    Description: "The unique ID of object from Terraform.",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// [C]: Create CSR
func resourceCreateCsr(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	var dnsnames []string
	var ips []net.IP
	var emails []string

	if err := d.Get("dnsnames").([]interface{}); len(err) > 0 {
		for _, dnsname := range d.Get("dnsnames").([]interface{}) {
			dnsnames = append(dnsnames, dnsname.(string))
		}
	} else {
		dnsnames = []string{}
	}

	if err := d.Get("ips").([]interface{}); len(err) > 0 {
		for _, ip := range d.Get("ips").([]interface{}) {
			ipaddr := net.ParseIP(ip.(string))
			ips = append(ips, ipaddr)
		}
	} else {
		ips = []net.IP{}
	}

	if err := d.Get("email").([]interface{}); len(err) > 0 {
		for _, email := range d.Get("email").([]interface{}) {
			emails = append(emails, email.(string))
		}
	} else {
		emails = []string{}
	}

	newsigner, err := NewDSMSigner(d.Get("kid").(string), dnsnames, ips, emails, d.Get("cn").(string), d.Get("ou").(string), d.Get("l").(string), d.Get("c").(string), d.Get("o").(string), d.Get("st").(string), d.Get("e").(string), m.(*api_client))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to create signer",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %v", newsigner),
		})
		return diags
	}

	generated_csr, err := newsigner.generate_csr()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to get DSM signer up",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %v", err),
		})
		return diags
	}

	if err := d.Set("value", generated_csr); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rand.Intn(99999999)))
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
