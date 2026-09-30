# Echo Scoreboard Visualizer

A small Windows scoreboard app for Echo VR. It serves a browser page at `http://localhost:8080` with the Blue and Orange teams, match score, round scores, clock, and up to four player slots per team.

The player panels show points (PTS), assists (AST), saves (SVS), stuns (STN), ping, and MVP. The first five values come from fields in the Echo session response. MVP is optional and may show `—` until a verified Spark data source is available. At zero time or when the match status indicates the post-game state, the center panel displays **Game Over**.

## Requirements

- Windows
- Go 1.27.1 or newer to run from source or build the executable
- Echo VR running for its local session API (`http://127.0.0.1:6721/session`)
- A browser, or OBS if using the page as a Browser Source

## Run From Source

Open PowerShell or Command Prompt in this project folder:

```powershell
go test ./...
go run ./cmd/scoreboard
```

Then open `http://localhost:8080` in a browser. Keep the terminal running while using the scoreboard; press `Ctrl+C` to stop it. The page refreshes its board data every five seconds. If Echo data is unavailable, it reports that it is waiting for a match.

## Build the Windows Executable

Run either build script from the project folder:

```powershell
.\build-scoreboard.ps1
```

Or, from Command Prompt:

```bat
build-scoreboard.bat
```

Both scripts create `builds\scoreboard.exe`. To launch it, run it from the project folder:

```powershell
.\builds\scoreboard.exe
```

The executable reads `web\index.html` at startup, so keep the `web` folder alongside the project and launch the executable with this project as its working directory.

## Use in OBS

1. Start the scoreboard app.
2. Add a **Browser** source to the scene.
3. Set its URL to `http://localhost:8080` and choose the desired width and height.

## Current Data Limitations

- The Echo client currently loads `testdata\echo_session.json` whenever it finds that file, before trying the live Echo endpoint. Since the sample is included in this repository, the current app will display the captured sample rather than live match data until fixture loading is limited to tests or otherwise disabled for normal runs.
- The app attempts an optional request to Spark at `http://localhost:6724/api`, but Spark's actual response has not been verified with a saved sample. MVP percentages may therefore be blank or unavailable.
- The server is intended for local use. Do not expose port 8080 to an untrusted network.
