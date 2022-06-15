FROM --platform=linux/amd64 hashicorp/terraform:latest
ARG FX_TF_VERSION

ENV FX_TF_VERSION=$FX_TF_VERSION
RUN mkdir -p /root/.terraform.d/plugins/fortanix.com/fortanix/dsm/${FX_TF_VERSION}/linux_amd64
COPY terraform-provider-dsm /root/.terraform.d/plugins/fortanix.com/fortanix/dsm/${FX_TF_VERSION}/linux_amd64
RUN chmod +x /root/.terraform.d/plugins/fortanix.com/fortanix/dsm/${FX_TF_VERSION}/linux_amd64/terraform-provider-dsm
