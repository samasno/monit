#!/bin/bash
cd $(dirname $0)
ls
out=$1
if [ -z $out ] 
then
    out=testing
fi;

targetDir=../certs
openssl genrsa -out $targetDir/$out.pem 2048
openssl req -new -x509 -key $targetDir/$out.pem -days 365 -subj "/C=/ST=/L=/O=/OU=/CN=localhost" -out $targetDir/$out.crt