@echo off
echo "protoc --go_out=. %1"
protoc --go_out=. %1