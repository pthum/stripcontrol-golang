#!/bin/sh
go test -coverprofile=./cov.out -json ./... 2>&1 | tee | go-junit-report -parser gojson > report.xml \
  && gocover-cobertura < cov.out > coverage.xml

