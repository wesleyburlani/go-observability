FROM golang:1.21

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt-get update && \
    apt-get install build-essential protobuf-compiler librdkafka-dev -y && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install github.com/spf13/cobra-cli@latest && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    wget https://github.com/ktr0731/evans/releases/download/v0.10.11/evans_linux_amd64.tar.gz && \
    tar -xzvf evans_linux_amd64.tar.gz && \
    mv evans ../bin && rm -f evans_linux_amd64.tar.gz

CMD ["tail", "-f", "/dev/null"]
