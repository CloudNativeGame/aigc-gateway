# syntax=docker/dockerfile:1
FROM golang:1.19 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -o aigc-gateway main.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose

FROM alpine:3.17

RUN apk add --no-cache ca-certificates bash expat curl \
  && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /app/aigc-gateway .
#COPY ./aigc-gateway .
COPY ./aigc-dashboard ./aigc-dashboard

# Run
CMD ["./aigc-gateway"]