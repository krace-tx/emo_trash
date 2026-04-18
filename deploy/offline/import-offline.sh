#!/usr/bin/env bash
# docker load restores images from docker save tarballs. (docker import is for docker export, different format.)
cd "$(dirname "$0")"
set +e
fail=0
docker load -i images/redis-7.2-alpine.tar || fail=1
docker load -i images/mysql-8.tar || fail=1
docker load -i images/mongo-6.0.tar || fail=1
docker load -i images/bitnamilegacy-etcd-3.5.18.tar || fail=1
docker load -i images/alpine-latest.tar || fail=1
docker load -i images/make-go-dev-latest.tar || fail=1
echo
read -n 1 -s -r -p "Press any key to close..."
echo
exit "$fail"
