FROM golang:1.12.0-alpine3.9 as builder
LABEL maintainer="Niroop Reddy <niroopreddy@outlook.com>"

RUN apk add git
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/niroopreddym/muxapiwithjwt

######## Start a new stage from scratch #######

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/muxapiwithjwt .

EXPOSE 9293

# Command to run the executable
CMD ["./muxapiwithjwt", "start"] 