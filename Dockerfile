FROM --platform=$BUILDPLATFORM golang:1.27 AS build
WORKDIR /go/src/app

COPY go.* .
RUN go mod download

COPY . .

RUN go vet -v ./...
RUN go test -v ./...

ARG TARGETOS TARGETARCH
ARG VERSION
RUN \
    VERSION=${VERSION:-`git tag --sort=-version:refname | head -n 1`} && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath \
    -ldflags "-s -w -X main.version=$VERSION" \
    -o HellPot ./cmd/HellPot


FROM gcr.io/distroless/static-debian13
LABEL org.opencontainers.image.source=https://github.com/SirTerrific/HellPot

COPY --from=build /go/src/app/HellPot /app
COPY --from=build /go/src/app/docker_config.toml /config
EXPOSE 8080
ENTRYPOINT ["/app", "-c", "/config"]
