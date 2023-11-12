#!/bin/bash

openssl version > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo "openssl was not found" >&2
    exit $?
fi

CERT_FILE="certificate.crt"
CSR_FILE="csr.csr"
KEY_FILE="private.key"

openssl genpkey -algorithm RSA -out "$KEY_FILE"
openssl req -new -key "$KEY_FILE" -out "$CSR_FILE"
openssl x509 -req -in "$CSR_FILE" -signkey "$KEY_FILE" -out "$CERT_FILE"
