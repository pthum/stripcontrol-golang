#!/bin/sh
set -x
runtests() {
  go test -coverprofile=./cov.out -json ./... 2>&1 | tee | go-junit-report -parser gojson > report.xml \
    && gocover-cobertura < cov.out > coverage.xml
}

runbuild() {
  GOOS=linux GOARCH=arm64 go build -o out/stripcontrol-app-aarch64 cmd/service/main.go
  GOOS=linux GOARCH=arm GOARM=7 go build -o out/stripcontrol-app-armv7 cmd/service/main.go
}

"$@"