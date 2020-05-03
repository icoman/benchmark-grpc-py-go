@echo off

SET HOST=localhost:6001
SET BS=2048
SET N=5

title Test go server insecured - %HOST%

cd go-server
start "go gRPC server insecure on %HOST%" go run . -host %HOST%
cd ..\go-client
go run . -host %HOST% -bs %BS% -n %N%
cd ..\py-client-server
echo.
python client.py -host %HOST% -bs %BS% -n %N%
pause
 