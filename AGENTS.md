# Echo scoreboard: agent instructions

A small scoreboard for Echo VR that shows the live score, player usernames, and MVP-style
stats, refreshed every 5 seconds. Built for its owner and one friend, not for distribution.


## Agent operating rule

- Always give a single-paragraph plan before making changes.
- Wait for an explicit "execute" from the user before editing files or running project work.
- If the environment is not ready, state the block honestly and wait for the user to unblock it.
- Do not skip the plan/execute step just because a prerequisite is now installed.

## What this is

- A single Windows `.exe` that polls Echo every 5 seconds and displays a native desktop
  scoreboard window. It does not host a webpage or listen on localhost.
- Written in **Go** and rendered with the installed WebView2 Runtime.
- Read-only. It sends GET requests only to Echo's local session API.

## What this is NOT

- Not a cheat or an overlay injected into the game. It's a separate window or web page.
- Not a web server. There is no app HTTP listener; Echo's API remains at
  `127.0.0.1:6721`.
- Not a replay tool. It works on live data only.

## Data sources

**Verify every field against a real response before you write code that uses it.** Don't
build from memory or guesses. Save one real response of each kind to `testdata/`, and
write parsing tests against those files.

1. **Echo VR API** (known): `GET http://127.0.0.1:6721/session`. It only responds while
   Echo is running and the player is in a match or spectating. It returns JSON with the
   game state, teams, players, and the score. It returns an error or refuses the connection
   when not in a match, so treat that as "waiting for match", not as a crash.
2. **Spark** is not used by the application. Its previously probed `/` and `/api` routes
  returned HTTP 404; do not assume an unverified Spark interface.
3. **MVP** (confirmed by owner, 2026-09-30): for `echo_arena`, calculate each player's
  score by summing `stats[award.stat] * award.multiplier` from the `awards` list in Echo's
  `r14/multiplayer/player_rewards.json`. The highest score is MVP. Do not display award
  titles in the app. Missing reward stats make the rating partial; don't invent values.

## Layout

```
cmd/scoreboard/main.go   # starts Echo poller + native desktop window
internal/echo/           # Echo API client + types (from testdata samples)
internal/board/          # board model and Echo Arena MVP calculation
internal/snapshot/       # local scoreboard PNG capture
web/index.html           # embedded desktop UI
testdata/                # saved Echo response used by parser tests
```

## Rules

- Poll Echo every 5 seconds, not faster.
- Every displayed number traces to a field in a saved `testdata/` sample. If a stat can't
  be traced, mark MVP as partial rather than inventing values.
- Save one scoreboard PNG per round or match-end event under the project `screenshots/`
  folder, with a fixed EST timestamp in its filename.
- Keep it small. No database, no accounts, no cloud.

## Commands

```
go test ./...                                                  # parsing tests against testdata/
go run ./cmd/scoreboard                                        # launch desktop window
GOOS=windows GOARCH=amd64 go build -o builds/scoreboard.exe ./cmd/scoreboard   # the .exe
```

## Done means

- With Echo running in a match: the desktop window shows the score, team players, stats,
  and MVP rating, updating within 5 seconds.
- With Echo closed: the window says "waiting for match" and doesn't crash.
- The app does not create an HTTP server or listen on port 8080.
- `go test ./...` passes against real saved samples.
