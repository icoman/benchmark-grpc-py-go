@echo off

SET protoc=D:\MyData\Kits\protoc-3.11.4-win64\bin\protoc

%protoc% -I . *.proto --go_out=plugins=grpc:.


pause

