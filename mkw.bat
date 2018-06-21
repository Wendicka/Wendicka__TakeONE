@echo off
cd cli

echo Compiling Builder
go build -o "../bin/wendicka_build.exe" -v wendicka_build.go


echo Compiling Quick Runtime Tool
go build -o "../bin/wendicka_run.exe" -v wendicka_run.go


cd ..