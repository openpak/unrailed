# Deploy: unrailed

## Running it

```
go run ./cmd/unrailedd                       # plain HTTP on :23500, health on :23501
go run ./cmd/unrailedd -tls-addr :8443 -cert chain.pem -key key.pem
```

| Variable | Meaning |
| --- | --- |
| `NX_INTERNAL_URL` | nx-baas internal API (default `http://127.0.0.1:20070`), used to check id tokens |
| `NX_INTERNAL_KEY` | its `X-Internal-Key` |

Ports are in [ports.md](ports.md). Logs never carry a token, account id or name, only
per-process tags (`redact`).

## Container

The Dockerfile builds a static `unrailedd` (`CGO_ENABLED=0`, `-trimpath`) into a
`gcr.io/distroless/static` image that runs as `nonroot:nonroot`, with `EXPOSE 23500 23501`.
