#!/bin/bash
export PG_NAME=fmdu
go build -ldflags="-s -w" -o $PG_NAME main.go
# バイナリの圧縮
strip $PG_NAME
upx -9 $PG_NAME

cp $PG_NAME /srv/www/fmdu/bin/.
