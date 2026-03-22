FROM golang:1.25-alpine AS builder

ENV PROTOC_VERSION="3.17.3"
ENV PROTOC_GEN_GO_VERSION="v1.36.11"

RUN apk add --no-cache curl git unzip tar

RUN curl -sfL -o /tmp/protoc.zip \
      "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" && \
    mkdir -p /opt/protoc && \
    unzip -q -d /opt/protoc /tmp/protoc.zip && \
    rm /tmp/protoc.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /opt/protoc/include/google /usr/local/include/google
COPY --from=builder /opt/protoc/bin/protoc /usr/local/bin/protoc
COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/protoc-gen-go

ENTRYPOINT ["/usr/local/bin/protoc"]