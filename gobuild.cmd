@echo off

set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0

echo Get version...
if not exist build.ver (
    echo Not found build.ver
    goto :EOF
)
set /p version=< build.ver
rem set yyyymmdd=%DATE:/=%
rem set version=%version% (%yyyymmdd%)
set version=%version%
echo Version : %version%
if ["%version%"]==[""] (
    echo Please check build.ver first.
    goto :EOF
)
echo Delete old module...
del agentserver.exe > nul 2>&1

echo golang mod init...
if not exist go.mod (
    go mod init
)
echo Build ...
go build -ldflags "-s -w -extldflags -static -X 'main.version=%version%'" -a .
echo Done.
