#!/bin/bash

if [ ! -e go.mod ]; then
go mod init
fi

if [ -e build.ver ]; then
read -d $'\x04' version < build.ver
echo $version
fi

echo Build ...
go build -ldflags "-s -w -extldflags -static -X 'main.version=$version'" -a .
