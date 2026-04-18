@echo off
rem One tar per image (docker save). Restore with docker load -i (not docker import).
cd /d "%~dp0"
if not exist images mkdir images
docker save -o images\redis-7.2-alpine.tar redis:7.2-alpine
docker save -o images\mysql-8.tar mysql:8
docker save -o images\mongo-6.0.tar mongo:6.0
docker save -o images\bitnamilegacy-etcd-3.5.18.tar bitnamilegacy/etcd:3.5.18
docker save -o images\alpine-latest.tar alpine:latest
docker save -o images\make-go-dev-latest.tar make-go-dev:latest
pause
