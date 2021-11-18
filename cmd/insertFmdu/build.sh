#!/bin/bash
export PG_NAME=insertFmdu
go build -ldflags="-s -w" -o $PG_NAME main.go
# バイナリの圧縮
strip $PG_NAME
upx -9 $PG_NAME
