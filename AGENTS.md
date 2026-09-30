# Echo scoreboard: agent instructions

A small scoreboard for Echo VR that shows the live score, player usernames, and MVP-style
stats, refreshed every 5 seconds. Built for its owner and one friend, not for distribution.

## Agent operating rule

- Always give a single-paragraph plan before making changes.
- Wait for an explicit "execute" from the user before editing files or running project work.
- If the environment is not ready, state the block honestly and wait for the user to unblock it.
- Do not skip the plan/execute step just because a prerequisite is now installed.

## What this is

- A single Windows `.exe` that polls local data sources every 5 seconds and serves a
  scoreboard page on `http://localhost:8080`. You open that page in a browser window, or
  add it as an OBS Browser Source.
- Written in **Go**, because Go builds one self-contained `.exe` with nothing to install.
- Read-only. It never sends anything to Echo or Spark except GET requests.

## What this is NOT

- Not a cheat or an overlay injected into the game. It's a separate window or web page.
- Not a public service. It only listens on `localhost`.
- Not a replay tool. It works on live data only.

## Data sources

**Verify every field against a real response before you write code that uses it.** Don't
build from memory or guesses. Save one real response of each kind to `testdata/`, and
write parsing tests against those files.

1. **Echo VR API** (known): `GET http://127.0.0.1:6721/session`. It only responds while
   Echo is running and the player is in a match or spectating. It returns JSON with the
   game state, teams, players, and the score. It returns an error or refuses the connection
   when not in a match, so treat that as "waiting for match", not as a crash.
2. **Spark** (to verify first): Spark runs a local server on `http://localhost:6724`
   (`/api/...` paths). The owner thinks the scoreboard data may come over a websocket, but
   that isn't confirmed. **Before any Spark code:** find out what Spark actually exposes
   (its docs, or by probing `localhost:6724` while Spark is running), save real samples to
   `testdata/`, and write down in this file what exists. Stats like MVP percentages may come
   only from Spark, or may have to be computed from the Echo API. Decide that from the
   samples, not before.
3. **MVP** (known, from dad, 2026-09-30): the in-match MVP is **the player with the highest
   point total of awards**. The awards and their point values are defined in Echo's
   `r14/multiplayer/player_rewards.json`. Up to 3 awards are displayed, plus MVP, and there
   is at least one unused award. Compute each player's award points live from the per-player
   `stats` block in `/session` (points, goals, assists, saves, stuns, steals, passes,
   catches, blocks, interceptions, shots_taken, possession_time), using the thresholds and
   points from that file. Don't make up weights.

## Layout

```
cmd/scoreboard/main.go   # starts the poller + web server
internal/echo/           # Echo API client + types (from testdata samples)
internal/spark/          # Spark client + types (only after the interface is verified)
internal/board/          # merges both sources into one Scoreboard struct
web/index.html           # the page; fetches /board.json every 5 s
testdata/                # real saved responses: echo_session.json, spark_*.json
```

## Rules

- Poll every 5 seconds, not faster. If a source is down, keep showing the last good data
  with a small "stale since HH:MM:SS" note. Don't blank the board.
- One source failing never breaks the other. Echo up + Spark down still shows the score.
- Every displayed number traces to a field in a saved `testdata/` sample. If a stat can't
  be traced, it doesn't go on the board.
- Keep it small. No database, no accounts, no cloud.

## Commands

```
go test ./...                                                  # parsing tests against testdata/
go run ./cmd/scoreboard                                        # run locally
GOOS=windows GOARCH=amd64 go build -o builds/scoreboard.exe ./cmd/scoreboard   # the .exe
```

## Done means

- With Echo running in a match: the score, both teams' usernames, and at least one stat
  per player show at `localhost:8080` and update within 5 seconds of a change.
- With Echo closed: the page says "waiting for match" and doesn't crash.
- `go test ./...` passes against real saved samples.
