@echo off
setlocal
cd /d "%dp0"
if not exist builds mkdir builds
go build -o .\builds\scoreboard.exe .\cmd\scoreboard
if errorlevel 1 (
    echo Build failed.
    exit /b 1
)
echo Build complete: %CD%\builds\scoreboard.exe
