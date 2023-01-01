## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /external-dns-docker

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /external-dns-docker /external-dns-docker

USER nonroot:nonroot

ENTRYPOINT ["/external-dns-docker"]
