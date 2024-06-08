FROM golang:alpine as builder

WORKDIR /build

# Add SSL certs to be transfered
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -v -o main

FROM scratch
WORKDIR /app
COPY --from=builder /build/main ./main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./main"]
