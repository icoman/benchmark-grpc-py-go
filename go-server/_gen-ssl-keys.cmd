@echo off

rem
rem https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
rem

SET PRGFOLDER=D:\MyData\Kits\SSL\_openssl\

rem generate some simple .key/.crt pairs 
%PRGFOLDER%\openssl genrsa -out server.key 2048
%PRGFOLDER%\openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -config %PRGFOLDER%\openssl.txt 


rem generate a certificate signing request (.csr) 
%PRGFOLDER%\openssl req -new -sha256 -key server.key -out server.csr -config %PRGFOLDER%\openssl.txt 
%PRGFOLDER%\openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 3650


pause
