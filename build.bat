@echo off
set DOCKER_BUILDKIT=1
if exist bin rmdir /S /Q bin
docker build --target bin --output bin/ .
