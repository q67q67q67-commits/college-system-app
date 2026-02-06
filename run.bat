@echo off
set "PATH=C:\Program Files\Go\bin;%PATH%"
cd /d "%~dp0"
go run ./cmd/api
pause
