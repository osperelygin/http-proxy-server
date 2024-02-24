#!/bin/sh
openssl req -new -key certs/cert.key -subj "/CN=$1" -sha256 | openssl x509 -req -days 365 -extfile <(printf "subjectAltName=DNS:$1") -CA certs/ca.crt -CAkey certs/ca.key -set_serial "0x`openssl rand -hex 8`" > certs/nck.crt
