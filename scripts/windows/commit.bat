@echo off
setLocal EnableDelayedExpansion

echo Staging changes
git add .

set /p message="Enter commit message: "

git commit -m "%message%"

echo Pushing changes
git push
cls