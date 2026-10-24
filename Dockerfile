# syntax=docker/dockerfile:1
# Unrailed!'s own web API (unrailed-online.com) for the Switch build. Plain HTTP on 23500;
# TLS for the client's :8443 is terminated in front, or by the service itself with -tls-addr.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOWORK=off go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/unrailedd ./cmd/unrailedd

FROM gcr.io/distroless/static
COPY --from=build /out/unrailedd /unrailedd
USER nonroot:nonroot
EXPOSE 23500 23501
ENTRYPOINT ["/unrailedd"]
