package dsm

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type dsmsigner struct {
	crypto.Signer
	crypto.Decrypter

	kid        string
	api_client *api_client

	Cn       string
	Ou       string
	L        string
	C        string
	O        string
	St       string
	Email    []string
	Dnsnames []string
	Ips      []net.IP
}

func (dsmsigner *dsmsigner) setProperty(propName string, propValue string) *dsmsigner {
	reflect.ValueOf(dsmsigner).Elem().FieldByName(propName).Set(reflect.ValueOf(propValue))
	return dsmsigner
}

func NewDSMSigner(kid string, dnsnames []string, ips []net.IP, email []string, cn string, ou string, l string, c string, o string, st string, api_client *api_client) (*dsmsigner, diag.Diagnostics) {
	var diags diag.Diagnostics

	var new_signer = &dsmsigner{
		kid:        kid,
		api_client: api_client,
		Dnsnames:   dnsnames,
		Ips:        ips,
		Email:      email,
	}

	if len(cn) > 0 {
		new_signer.setProperty("Cn", cn)
	}
	if len(ou) > 0 {
		new_signer.setProperty("Ou", ou)
	}
	if len(l) > 0 {
		new_signer.setProperty("L", l)
	}
	if len(c) > 0 {
		new_signer.setProperty("C", c)
	}
	if len(st) > 0 {
		new_signer.setProperty("St", st)
	}
	if len(o) > 0 {
		new_signer.setProperty("O", o)
	}

	return new_signer, diags
}

// generate_csr: Generate CSR locally
func (dsmsigner *dsmsigner) generate_csr() (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var subj pkix.Name
	subj.CommonName = dsmsigner.Cn

	if dsmsigner.C != "" {
		subj.Country = []string{dsmsigner.C}
	}
	if dsmsigner.St != "" {
		subj.Province = []string{dsmsigner.St}
	}
	if dsmsigner.L != "" {
		subj.Locality = []string{dsmsigner.L}
	}
	if dsmsigner.O != "" {
		subj.Organization = []string{dsmsigner.O}
	}
	if dsmsigner.Ou != "" {
		subj.OrganizationalUnit = []string{dsmsigner.Ou}
	}

	rawSubj := subj.ToRDNSequence()

	asn1Subj, _ := asn1.Marshal(rawSubj)
	// FYOO: SAN support
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		DNSNames:           dsmsigner.Dnsnames,
		EmailAddresses:     dsmsigner.Email,
		IPAddresses:        dsmsigner.Ips,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, dsmsigner)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to generate CSR",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %s", err),
		})
		return "", diags
	}
	generated_csr := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	generated_csr_pem := string(pem.EncodeToMemory(&generated_csr))

	return generated_csr_pem, diags
}

func (dsmsigner dsmsigner) Public() crypto.PublicKey {
	req, _, err := dsmsigner.api_client.APICall("GET", fmt.Sprintf("crypto/v1/keys/%s", dsmsigner.kid))
	if err != nil {
		panic("Unable to call DSM")
	}

	fxPubKeyDer, _ := base64.StdEncoding.DecodeString(req["pub_key"].(string))

	fxPubKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: fxPubKeyDer,
	}

	fxPubKeyPem := string(pem.EncodeToMemory(&fxPubKeyBlock))

	pubKeyBlock, _ := pem.Decode([]byte(fxPubKeyPem))
	if pubKeyBlock == nil {
		panic("failed to parse PEM block containing the public key")
	}

	pubKey, err1 := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err1 != nil {
		panic("failed to parse DER encoded public key: " + err1.Error())
	}

	return pubKey
}

func (dsmsigner dsmsigner) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	hash := opts.HashFunc()
	if len(digest) != hash.Size() {
		return nil, fmt.Errorf("[DSM SDK]: signer: Digest length doesn't match crypto algorithm")
	}

	sign_op := map[string]interface{}{
		"kid":      dsmsigner.kid,
		"hash_alg": "SHA256",
		"data":     base64.StdEncoding.EncodeToString(digest),
	}

	reqfpi, err := dsmsigner.api_client.FindPluginId("Terraform Plugin - CSR")
	if err != nil {
		return nil, fmt.Errorf("[DSM SDK]: signer: Unable to call DSM provider API client: GET: sys/v1/plugins: %v", err)
	}
	var endpoint = fmt.Sprintf("sys/v1/plugins/%s", string(reqfpi))
	var operation = "POST"

	req, err := dsmsigner.api_client.APICallBody(operation, endpoint, sign_op)
	if err != nil {
		return nil, fmt.Errorf("[DSM SDK]: signer: Unable to call DSM provider API client: POST: sys/v1/plugins: %v", err)
	}

	signature, err1 := base64.StdEncoding.DecodeString(req["signature"].(string))
	if err1 != nil {
		return nil, fmt.Errorf("[DSM SDK]: unable to convert from base64")
	}
	return signature, nil
}
