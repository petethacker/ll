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

if exist c:\bin\python\modules\sign.py (
    c:\bin\python\3.7.4\python.exe c:\bin\python\modules\sign.py c:\bin\apps\ll_arm64.exe
)
if %errorlevel% neq 0 (
    popd
    exit /b %errorlevel%
)

:Success
popd
exit /b 0