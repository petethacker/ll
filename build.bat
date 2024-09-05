@echo off

set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

pushd "%~dp0"

go build -o c:\bin\apps\ll.exe -ldflags="-s -w"
if %errorlevel% neq 0 (
    popd
    exit /b %errorlevel%
)

:: if exist c:\bin\python\modules\sign.py (
::    c:\bin\python\3.10.7\python.exe c:\bin\python\modules\sign.py c:\bin\apps\ll.exe
:: )
:: if %errorlevel% neq 0 (
::     popd
::     exit /b %errorlevel%
:: )

call build_arm64.bat
if %errorlevel% neq 0 (
    popd
    exit /b %errorlevel%
)

popd
exit /b 0