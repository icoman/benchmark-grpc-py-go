@echo off

SET HOST=localhost:6003
SET BS=2048
SET N=5

title Test py server insecured - %HOST%

cd py-client-server
start "py gRPC server insecure on %HOST%" python server.py -host %HOST%
cd ..\go-client
go run . -host %HOST% -bs %BS% -n %N%
cd ..\py-client-server
echo.
python client.py -host %HOST% -bs %BS% -n %N%
pause
 