# Changelog — unrailed

## v0.1.0 — 2026-11-02

- Release workflow publishing to ghcr.io; replaces the hand-built production image.

- Add a release workflow publishing the image to ghcr.io on v*.*.* tags.


- `unrailedd`: the client's web API on unrailed-online.com. The online gate checks the player's
  id token against nx-baas, the public-match block list is empty, kick reports and avatar/replay
  uploads are accepted and not stored. Optional TLS listener for the client's `:8443`.
