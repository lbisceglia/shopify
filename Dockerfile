# STEP 1 : Build the executable binary
# ================================
FROM golang:1.17-alpine3.15 as builder
LABEL stage=builder

# All these steps will be cached
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .

# Get dependencies; also cached unless go.mod/go.sum change.
RUN go mod download
RUN go mod verify

# COPY the source code
COPY . .

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/shopify

# # STEP 2: Build a small image
# # ===========================
# FROM scratch as shopify

# # Copy static executable
# COPY --from=builder /go/bin/shopify /go/bin/shopify
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run the binary
ENTRYPOINT ["/go/bin/shopify"]
