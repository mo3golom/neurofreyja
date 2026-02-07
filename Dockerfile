# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.25.7-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w" -o /out/bot ./cmd/bot

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
RUN adduser -D -u 10001 app

COPY --from=build /out/bot /usr/local/bin/bot

USER app
ENTRYPOINT ["/usr/local/bin/bot"]
