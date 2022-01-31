@echo off
echo "protoc --go_out=plugins=grpc:. %1"
protoc --go_out=plugins=grpc:. %1