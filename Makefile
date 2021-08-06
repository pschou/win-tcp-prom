PROG_NAME := "win-tcp-prom"
VERSION = 0.1.$(shell date +%Y%m%d.%H%M)
SRC = $(shell git config --get remote.origin.url | sed 's/.*@//;s/:/\//;s/\.git//' )
FLAGS := "-s -w -X main.version=${VERSION}@${SRC}"

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags=${FLAGS} -o ${PROG_NAME}.exe
#	upx --lzma ${PROG_NAME}.exe

