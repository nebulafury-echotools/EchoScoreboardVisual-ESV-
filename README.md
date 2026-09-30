# Echo Scoreboard Visualizer

A standalone Windows scoreboard for Echo VR. Launch the executable to open its desktop window; it does not host a webpage or listen on a localhost port. The embedded UI shows Blue and Orange teams, match score, rounds won, timer, and up to four players per team.

Player panels show PTS, AST, SVS, STN, PING, and the calculated MVP score. Individual award titles are not displayed. For `Echo_Arena`, the MVP score is the sum of each award stat multiplied by its configured multiplier from Echo's `r14/multiplayer/player_rewards.json`. The highest score is marked MVP. A `~` prefix means one or more required stats were absent from `/session`, so the rating is partial. Outside `Echo_Arena`, MVP displays as `N/A`.

At round or match end, the app saves one PNG of its scoreboard window for that end event to `%LOCALAPPDATA%\EchoScoreboardVisual\screenshots`.

## Requirements

- Windows with the Microsoft Edge WebView2 Runtime installed
- Go 1.27.1 or newer to run from source or build the executable
- Echo VR with its local session API enabled at `http://127.0.0.1:6721/session`

## Run From Source

From PowerShell or Command Prompt in this project folder:

```powershell
go test ./...
go run ./cmd/scoreboard
```

The scoreboard opens in a desktop window and refreshes from Echo every five seconds. Echo must be in a match or spectating; otherwise the window says "Waiting for match." Close the window to exit.

## Build the Windows Executable

Run either build script from the project folder:

```powershell
.\build-scoreboard.ps1
```

Or from Command Prompt:

```bat
build-scoreboard.bat
```

Both scripts create `builds\scoreboard.exe`. Launch it with:

```powershell
.\builds\scoreboard.exe
```

The UI is embedded in the executable. The WebView2 Runtime is a Windows prerequisite and is not bundled into the EXE.

## OBS

Use a **Window Capture** source and select **Echo Scoreboard**. No Browser Source or localhost URL is needed.

## MVP Data Notes

The saved Echo session fixture contains only a subset of the `echo_arena` reward stats, so fixture-based ratings are partial. At runtime, the app polls Echo directly and does not query Spark. Spark's previously probed `/` and `/api` routes returned HTTP 404; no Spark endpoint is used for MVP.
