#!/bin/sh

set -ex

DOMAIN1='tictac.localhost'
DOMAIN2='add.localhost'

DIR=$(dirname $0)

# Generate the CA Key and Certificate:
# 
openssl req -x509 -sha256 -newkey rsa:4096 -keyout "${DIR}/ca.key" -out "${DIR}/ca.crt" -days 3650 -nodes -subj '/CN=LOCALHOST Cert Authority'

# Generate the Server Key, and Certificate and Sign with the CA Certificate:
# 
openssl req -new -newkey rsa:4096 -keyout "${DIR}/server1.key" -out "${DIR}/server1.csr" -nodes -subj '/CN='"${DOMAIN1}"
openssl x509 -req -sha256 -days 365 -in "${DIR}/server1.csr" -CA "${DIR}/ca.crt" -CAkey "${DIR}/ca.key" -set_serial 01 -out "${DIR}/server1.crt"

cat "${DIR}/server1.crt" "${DIR}/ca.crt" > "${DIR}/${DOMAIN1}.chain.crt"

openssl req -new -newkey rsa:4096 -keyout "${DIR}/server2.key" -out "${DIR}/server2.csr" -nodes -subj '/CN='"${DOMAIN2}"
openssl x509 -req -sha256 -days 365 -in "${DIR}/server2.csr" -CA "${DIR}/ca.crt" -CAkey "${DIR}/ca.key" -set_serial 01 -out "${DIR}/server2.crt"

cat "${DIR}/server2.crt" "${DIR}/ca.crt" > "${DIR}/${DOMAIN2}.chain.crt"