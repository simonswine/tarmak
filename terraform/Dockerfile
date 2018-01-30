# Copyright Jetstack Ltd. See LICENSE for details.
FROM alpine:3.6

RUN apk add --no-cache unzip curl

# install terraform
ENV TERRAFORM_VERSION 0.10.8
ENV TERRAFORM_HASH b786c0cf936e24145fad632efd0fe48c831558cc9e43c071fffd93f35e3150db
RUN curl -sL  https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip > /tmp/terraform.zip && \
    echo "${TERRAFORM_HASH}  /tmp/terraform.zip" | sha256sum  -c && \
    unzip /tmp/terraform.zip && \
    rm /tmp/terraform.zip && \
    mv terraform /usr/local/bin/terraform && \
    chmod +x /usr/local/bin/terraform

# install terraform plugins
WORKDIR /terraform

ENV TERRAFORM_PROVIDER_AWSTAG_VERSION 0.1
ENV TERRAFORM_PROVIDER_AWSTAG_HASH f2f1a153637f837ab521a4f1d0e236e08927354610495a199e4f67985de28a83
RUN curl -sL https://github.com/dippynark/terraform-provider-awstag/releases/download/release-${TERRAFORM_PROVIDER_AWSTAG_VERSION}/terraform-provider-awstag_linux_amd64.zip > /tmp/terraform-provider-awstag.zip && \
    sha256sum /tmp/terraform-provider-awstag.zip && \
    echo "${TERRAFORM_PROVIDER_AWSTAG_HASH}  /tmp/terraform-provider-awstag.zip" | sha256sum -c && \
    unzip /tmp/terraform-provider-awstag.zip && \
    rm /tmp/terraform-provider-awstag.zip && \
    mv terraform-provider-awstag /terraform/terraform-provider-awstag && \
    chmod +x /terraform/terraform-provider-awstag

ADD providers.tf /terraform
RUN terraform init
RUN rm providers.tf
