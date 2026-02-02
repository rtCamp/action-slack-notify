FROM golang:1.23-alpine3.19@sha256:5f3336882ad15d10ac1b59fbaba7cb84c35d4623774198b36ae60edeba45fd84 AS builder

LABEL "com.github.actions.icon"="bell"
LABEL "com.github.actions.color"="yellow"
LABEL "com.github.actions.name"="Slack Notify"
LABEL "com.github.actions.description"="This action will send notification to Slack"
LABEL "org.opencontainers.image.source"="https://github.com/rtCamp/action-slack-notify"

WORKDIR ${GOPATH}/src/github.com/rtcamp/action-slack-notify
COPY main.go ${GOPATH}/src/github.com/rtcamp/action-slack-notify

ENV CGO_ENABLED 0
ENV GOOS linux

RUN go build -a -installsuffix cgo -ldflags '-w  -extldflags "-static"' -o /go/bin/slack-notify main.go

# alpine:latest as of 2024-03-11
FROM alpine@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

COPY --from=builder /go/bin/slack-notify /usr/bin/slack-notify

ENV VAULT_VERSION 1.0.2

RUN apk update \
	&& apk upgrade \
	&& apk add \
		bash \
		jq \
		ca-certificates \
		python3 \
		py3-pip \
		rsync \
	&& python3 -m pip install --break-system-packages shyaml \
	&& rm -rf /var/cache/apk/*

# Setup Vault
RUN wget https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip && \
	unzip vault_${VAULT_VERSION}_linux_amd64.zip && \
	rm vault_${VAULT_VERSION}_linux_amd64.zip && \
	mv vault /usr/local/bin/vault

# fix the missing dependency - https://stackoverflow.com/a/35613430
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY *.sh /

RUN chmod +x /*.sh

ENTRYPOINT ["/entrypoint.sh"]
