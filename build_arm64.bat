@echo off

set CGO_ENABLED=1
set GOOS=windows
set GOARCH=arm64

pushd "%~dp0"

go build -o c:\bin\apps\ll_arm64.exe
if %errorlevel% neq 0 (
    popd
    exit /b %errorlevel%
)

popd
exit /b 0