FROM golang:1.14-alpine3.11@sha256:6578dc0c1bde86ccef90e23da3cdaa77fe9208d23c1bb31d942c8b663a519fa5 AS builder

LABEL "com.github.actions.icon"="bell"
LABEL "com.github.actions.color"="yellow"
LABEL "com.github.actions.name"="Slack Notify"
LABEL "com.github.actions.description"="This action will send notification to Slack"
LABEL "org.opencontainers.image.source"="https://github.com/rtCamp/action-slack-notify"

WORKDIR ${GOPATH}/src/github.com/rtcamp/action-slack-notify
COPY main.go ${GOPATH}/src/github.com/rtcamp/action-slack-notify

ENV CGO_ENABLED 0
ENV GOOS linux

RUN go get -v ./... && \
	go build -installsuffix cgo -ldflags '-w  -extldflags "-static"' -o /go/bin/slack-notify .

# alpine:latest at 2023-08-23T11:48:25.468623934UTC
FROM alpine@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a

ENV VAULT_VERSION 1.0.2

# Slow, rarely-changing layers first so they stay cached when code changes
RUN apk upgrade --no-cache \
	&& apk add --no-cache \
		bash \
		jq \
		ca-certificates \
		python3 \
		py3-pip \
		rsync \
	&& pip3 install --no-cache-dir shyaml \
	&& wget -q https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip \
	&& unzip vault_${VAULT_VERSION}_linux_amd64.zip -d /usr/local/bin vault \
	&& rm vault_${VAULT_VERSION}_linux_amd64.zip \
	&& mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Frequently-changing layers last
COPY --from=builder /go/bin/slack-notify /usr/bin/slack-notify
COPY *.sh /
RUN chmod +x /*.sh

ENTRYPOINT ["/entrypoint.sh"]
