@echo off

rem
rem https://grpc.io/docs/tutorials/basic/python/
rem pip install grpcio-tools
rem

python -m grpc_tools.protoc -I../srvapi --python_out=. --grpc_python_out=. ../srvapi/*.proto

pause
