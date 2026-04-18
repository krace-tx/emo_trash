@echo off
rem docker load restores images from docker save tarballs. (docker import is for docker export, different format.)
cd /d "%~dp0"
docker load -i images\redis-7.2-alpine.tar
docker load -i images\mysql-8.tar
docker load -i images\mongo-6.0.tar
docker load -i images\bitnamilegacy-etcd-3.5.18.tar
docker load -i images\alpine-latest.tar
docker load -i images\make-go-dev-latest.tar
pause
