FROM alpine:3.6

RUN apk add --no-cache unzip curl

# install packer
ENV PACKER_VERSION 1.1.0
ENV PACKER_HASH bd1eddfa114f7e6258ef3419613380297f1b4e438f5bae92f1177150519be934
RUN curl -sL  https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip > /tmp/packer.zip && \
    echo "${PACKER_HASH}  /tmp/packer.zip" | sha256sum  -c && \
    unzip /tmp/packer.zip && \
    rm /tmp/packer.zip && \
    mv packer /usr/local/bin/packer && \
    chmod +x /usr/local/bin/packer
