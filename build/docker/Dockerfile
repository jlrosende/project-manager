# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.4-alpine

FROM golang:${GO_VERSION} as build

ENV GOCACHE=/go-cache \
    GOMODCACHE=/gomod-cache \
    CGO_ENABLED=0

WORKDIR /src

ARG GOMODCACHE GOCACHE

RUN --mount=type=cache,target="$GOMODCACHE" \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETOS TARGETARCH 

RUN --mount=type=cache,target="$GOMODCACHE" \
    --mount=type=cache,target="$GOCACHE" \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -o /bin/pm ./cmd/main.go

FROM scratch

COPY --from=build /bin/pm /pm

CMD [ "/pm" ]

