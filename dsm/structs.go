// **********
// Terraform Provider - DSM: common: structs
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.0
//       - Date:      27/11/2020
// **********

package dsm

// [-] Structs to define DSM AWS Group
type AWSGroup struct {
	Acct_id        string                 `json:"acct_id"`
	Creator        DSMCreator             `json:"creator"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Group_id       string                 `json:"group_id"`
	Hmg            map[string]AWSGroupHmg `json:"hmg"`
	Hmg_redundancy string                 `json:"hmg_redundancy"`
}

type AWSGroupHmg struct {
	Kind       string `json:"kind"`
	Url        string `json:"url"`
	Access_key string `json:"access_key"`
	Hsm_order  int    `json:"hsm_order"`
}

// [-] Structs to define DSM AWS Security Object
type AWSSobject struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Key_size        int
	Acct_id         string
	Group_id        string
	Creator         DSMCreator
	Kid             string
	Obj_type        string
	Custom_metadata AWSSobjectCustomMetadata
	Enabled         bool
	State           string
	External        AWSSobjectExternal
	Links           DSMSobjectLinks
}

type AWSSobjectCustomMetadata struct {
	Aws_key_state     string `json:"aws-key-state"`
	Aws_aliases       string `json:"aws-aliases"`
	Aws_deletion_date string `json:"aws-deletion-date"`
}

type AWSSobjectExternal struct {
	Hsm_group_id string
	Id           AWSSobjectExternalId
}

type AWSSobjectExternalId struct {
	Key_arn string
	Key_id  string
}

// [-] Structs to define DSM Azure Group
type AzureGroup struct {
	Acct_id        string                   `json:"acct_id"`
	Creator        DSMCreator               `json:"creator"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Group_id       string                   `json:"group_id"`
	Hmg            map[string]AzureGroupHmg `json:"hmg"`
	Hmg_redundancy string                   `json:"hmg_redundancy"`
}

type AzureGroupHmg struct {
	Kind            string `json:"kind"`
	Url             string `json:"url"`
	Key_vault_type  string `json:"key_vault_type"`
	Client_id       string `json:"client_id"`
	Subscription_id string `json:"subscription_id"`
	Tenant_id       string `json:"tenant_id"`
	Hsm_order       int    `json:"hsm_order"`
}

// [-] Structs to define DSM Azure Security Object
type AzureSobject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key_size    int
	Acct_id     string
	Group_id    string
	Creator     DSMCreator
	Kid         string
	Obj_type    string
	//	Custom_metadata AWSSobjectCustomMetadata
	Enabled bool
	State   string
	//	External        AWSSobjectExternal
	Links DSMSobjectLinks
}

// [-] Structs to define DSM GCP Security Object
type GCPSobject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key_size    int
	Acct_id     string
	Group_id    string
	Creator     DSMCreator
	Kid         string
	Obj_type    string
	Enabled     bool
	State       string
	Links       DSMSobjectLinks
}

// [-] Structs to define DSM definition
type DSMCreator struct {
	User string
}

type DSMSobjectLinks struct {
	Copiedfrom string
}
