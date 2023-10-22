FROM debian AS sqlc-generate
COPY --from=kjconroy/sqlc /workspace/sqlc /usr/bin/sqlc
WORKDIR /app
COPY . .
RUN sqlc generate

FROM golang:1.21 as vendor
COPY --from=sqlc-generate /app /app
WORKDIR /app
RUN go get -d -v ./...

FROM vendor as builder
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch
WORKDIR /app/bin
COPY --from=builder /app/bin/api ./
EXPOSE 3000
ENTRYPOINT ["./api"]
