# unrailed — OpenPak's backend for Unrailed!

Unrailed! on Nintendo Switch (`010042C011476000`) plays online through two services:

| What | Where | Served by |
| --- | --- | --- |
| Matchmaking, lobbies, gameplay | Photon Realtime, `ns.exitgames.com` | [`photon`](https://github.com/openpak/photon) |
| Online gate, kick list, avatar and replay uploads | the developer's web API, `unrailed-online.com:8443` / `:8080` | this repo (`unrailedd`) |

The game is host-authoritative: the Photon master client runs the session and the others follow
it through room events or Photon Punchthrough. Nothing here simulates the game.

## Status

Not deployed and not yet tested against the game. The contract in [docs/protocol.md](docs/protocol.md)
was read from the Switch executable. It still needs a capture of a real session.

## Running it

```
go run ./cmd/unrailedd                       # plain HTTP on :23500, health on :23501
go run ./cmd/unrailedd -tls-addr :8443 -cert chain.pem -key key.pem
```

| Variable | Meaning |
| --- | --- |
| `NX_INTERNAL_URL` | nx-baas internal API (default `http://127.0.0.1:20070`), used to check id tokens |
| `NX_INTERNAL_KEY` | its `X-Internal-Key` |

Ports are in [docs/ports.md](docs/ports.md). Logs never carry a token, account id or name, only
per-process tags (`redact`).

## License

AGPL-3.0-only. See [LICENSE](LICENSE).
