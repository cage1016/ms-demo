#!/bin/sh

set -ex

DIR=$(dirname $0)

kubectl create secret tls tictac-tls-secret --key "${DIR}/server1.key" --cert "${DIR}/tictac.localhost.chain.crt"
kubectl create secret tls add-tls-secret --key "${DIR}/server2.key" --cert "${DIR}/add.localhost.chain.crt"