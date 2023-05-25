# Step 1: Modules caching
FROM golang:1.20.4-alpine3.18 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.20.4-alpine3.18 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

# Step 3: Final
FROM golang:1.20.4-alpine3.18
COPY --from=builder /bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app", "-n"]
CMD ["6"]