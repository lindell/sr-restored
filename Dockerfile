FROM golang:alpine as gobuilder

WORKDIR /build

# Add SSL certs to be transfered
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -v -o main

FROM node:20-alpine AS sveltebuilder

WORKDIR /app

ENV PUBLIC_BASE_URL="https://sr-uncensored.fly.dev"

COPY ./frontend/package*.json .
RUN npm ci
COPY ./frontend .
RUN npm run build
RUN npm prune --production

FROM scratch
WORKDIR /app

COPY --from=gobuilder /build/main ./main
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=sveltebuilder /app/build ./static

CMD ["./main"]
