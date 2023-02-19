#!/bin/sh
CGO_ENABLED=0 go test -coverprofile=./cov.out -json ./... 2>&1 | tee | go-junit-report -parser gojson > report.xml \
  && gocover-cobertura < cov.out > coverage.xml

