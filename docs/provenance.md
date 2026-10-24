# Provenance

Where each fact in [protocol.md](protocol.md) came from, per the
[clean-room policy](https://docs.openpak.org/contributors/clean-room/). Only the owner's own game
client was examined. No server code of the developer was seen or used.

| Fact | Source |
| --- | --- |
| Base URLs `https://unrailed-online.com:8443`, `http://unrailed-online.com:8080` and the setting that picks one | string literals and the getter that returns one of them, Switch executable v1.0.12 |
| Paths `/nintendoHasSubscription?ticket=`, `/getKick`, `/notifyKick?id={0}&name={1}`, `/uploadAvatar`, `/uploadReplay` | string literals, Switch executable |
| Query `userid=`, `&ticket=`, `&buildNumber=` appended to the kick calls | the query builder that references those literals, Switch executable |
| `NoOnlineService` is the only refusal body | the response check after the subscription request (a 15-character comparison), Switch executable |
| Kick list is `id:name` entries | the response parser after `/getKick` (split, then `:`), Switch executable |
| Multipart field names for the two uploads | literals referenced by the upload code, Switch executable |
| Photon AppId and name server | literals in the region/name-server setup, Switch executable |
| The live host is an HTTP server on both ports | one `GET /` to each port, 2026-10-24 (404 from both) |
