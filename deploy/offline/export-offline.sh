#!/usr/bin/env bash
# One tar per image (docker save). Restore with docker load -i (not docker import).
cd "$(dirname "$0")"
mkdir -p images
set +e
fail=0
docker save -o images/redis-7.2-alpine.tar redis:7.2-alpine || fail=1
docker save -o images/mysql-8.tar mysql:8 || fail=1
docker save -o images/mongo-6.0.tar mongo:6.0 || fail=1
docker save -o images/bitnamilegacy-etcd-3.5.18.tar bitnamilegacy/etcd:3.5.18 || fail=1
docker save -o images/alpine-latest.tar alpine:latest || fail=1
docker save -o images/make-go-dev-latest.tar make-go-dev:latest || fail=1
echo
read -n 1 -s -r -p "Press any key to close..."
echo
exit "$fail"
