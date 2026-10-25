# unrailed-online.com — the client's contract

Read from the Switch executable (Unrailed! `010042C011476000`, update v1.0.12, a .NET NativeAOT
build), see [provenance](provenance.md). Keep this file in step with `cmd/unrailedd`.

## Base URL

The client builds every request on one of two base URLs, chosen by a setting:

- `https://unrailed-online.com:8443`
- `http://unrailed-online.com:8080`

Both ports are named by the client, so the redirected name has to answer on both.

## Requests

Every query value except `buildNumber` names or authenticates a player and is never logged raw.
`ticket` is the player's Nintendo Account id token (the `nn::account` id token cache).

| Request | Client reads | unrailedd answers |
| --- | --- | --- |
| `GET /nintendoHasSubscription?ticket=<id token>` | body text; exactly `NoOnlineService` refuses online play, anything else (or no body) allows it | `OK`, or `NoOnlineService` when nx-baas rejects the token or the account is banned; an nx-baas outage answers `OK` |
| `GET /getKick?userid=<id>&ticket=<id token>&buildNumber=<n>` | lines of `id:name` — the players blocked from public matches | empty body: nobody is blocked |
| `GET /notifyKick?id=<id>&name=<name>&userid=…&ticket=…&buildNumber=…` | success only | `200`, logged with tags |
| `POST /uploadAvatar` multipart `ticket`, `userid`, `avatar` (JPEG) | success only | `200`; nothing stored |
| `POST /uploadReplay` multipart `userid`, `ticket`, `datane`, `replay`, `week` | success only | `200`; nothing stored (only sent when the player opted in) |

The run summary that rides with a replay is JSON with `build`, `score`, `bolts`, `hardcore`,
`players`, `difficulty`, `seed`, `seedString`, `mode`, `crashes`, `uniqueGameId`,
`screenshotPNG` and `maxPlayers`.

## What is not here

Matchmaking, lobbies and gameplay are Photon Realtime on `ns.exitgames.com` (global AppId
`9d45adf2-e83f-4281-b9f5-02728f813f80`), served by `platforms/switch/games/photon`. Quick match is
`OpJoinRandomGame` with a Photon SQL-lobby filter on room properties `C0`, `C1` and `C3`. Private
sessions join a room by its session id. High scores are kept on the console.

## From the previous README

Unrailed! on Nintendo Switch (`010042C011476000`) plays online through two services:

| What | Where | Served by |
| --- | --- | --- |
| Matchmaking, lobbies, gameplay | Photon Realtime, `ns.exitgames.com` | [`photon`](https://github.com/openpak/photon) |
| Online gate, kick list, avatar and replay uploads | the developer's web API, `unrailed-online.com:8443` / `:8080` | this repo (`unrailedd`) |

The game is host-authoritative: the Photon master client runs the session and the others follow
it through room events or Photon Punchthrough. Nothing here simulates the game.
