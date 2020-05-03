@echo off

SET HOST=localhost:6004
SET BS=2048
SET N=5

title Test py server secured - %HOST%

cd py-client-server
start "py gRPC server secure on %HOST%" python server.py -host %HOST% -secured
cd ..\go-client
go run . -host %HOST% -bs %BS% -n %N% -secured
cd ..\py-client-server
echo.
python client.py -host %HOST% -bs %BS% -n %N% -secured
pause
 