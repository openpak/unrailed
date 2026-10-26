# Unrailed!

OpenPak's Unrailed! Server.

| | |
| --- | --- |
| Game | [Unrailed!](https://openpak.org/titles#t-138) |
| Network | Nintendo Switch |
| Title ID | `010042C011476000` |
| Licence | AGPL-3.0-only |

## What it serves

| Host | What | Served by |
| --- | --- | --- |
| `ns.exitgames.com` | Matchmaking, lobbies, gameplay (Photon Realtime) | [`photon`](https://github.com/openpak/photon) |
| `unrailed-online.com` `:8443` / `:8080` | Online gate, kick list, avatar and replay uploads | this repo (`unrailedd`) |

## Run

```sh
go run ./cmd/unrailedd
```

| Variable | Meaning |
| --- | --- |
| `NX_INTERNAL_URL` | nx-baas internal API (default `http://127.0.0.1:20070`), used to check id tokens |
| `NX_INTERNAL_KEY` | its `X-Internal-Key` |

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)
- [Ports](docs/ports.md)
- [Provenance](docs/provenance.md)

Live compatibility status: [openpak.org/titles#t-138](https://openpak.org/titles#t-138).

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
