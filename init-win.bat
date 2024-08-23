@echo off
cd C:\Users\guido\Documents\heritage

REM Update repository
git pull

echo Repo updated, waiting 5 seconds because windows is a shit.

timeout /t 5 /nobreak

make run

pause
