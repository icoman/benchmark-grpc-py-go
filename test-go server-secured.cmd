@echo off

SET HOST=localhost:6002
SET BS=2048
SET N=5

title Test go server secured - %HOST%

cd go-server
start "go gRPC server secure on %HOST%" go run . -host %HOST% -secured
cd ..\go-client
go run . -host %HOST% -bs %BS% -n %N% -secured
cd ..\py-client-server
echo.
python client.py -host %HOST% -bs %BS% -n %N% -secured
pause
 