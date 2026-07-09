@echo off
setLocal EnableDelayedExpansion

echo Removing all docker containers and volumes
for /f "delims=" %%i in ('docker ps -aq') do docker rm -f %%i
for /f "delims=" %%i in ('docker volume ls -q') do docker volume rm %%i

set /p choice="Do you want to start containers? (y/n) "

if /i "%choice%" neq "y" (
    echo container initialization cancelled
) else (
    echo Starting docker compose
    docker compose up -d
)