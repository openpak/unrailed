# Ports

| Port | Use |
| --- | --- |
| 23500 | plain HTTP API (the client's `:8080`, and what a TLS terminator for `:8443` forwards to) |
| 23501 | `/health` |
| 23502–23509 | reserved |
| 8443/tcp, 8080/tcp | fixed by the client (`unrailed-online.com:8443` TLS, `:8080` plain); published on the host and forwarded to 23500 (or `-tls-addr :8443` when the service terminates TLS) |

23500–23509 is the first block of the 23500–23999 range; the third-party block 22600–22699 is full.
