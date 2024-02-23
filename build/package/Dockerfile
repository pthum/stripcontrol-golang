FROM golang:1.21-alpine3.19 AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/stripcontrol-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download && go mod tidy

COPY . .

FROM build_base AS test
# switch to head instead of latest release due to https://github.com/jstemmer/go-junit-report/issues/138
RUN go install github.com/jstemmer/go-junit-report/v2@bfac3ec01f0cf6ed0dbff9f5f2cb48eb04db38a7 \
   && go install github.com/t-yuki/gocover-cobertura@latest \
# Unit tests
  && chmod +x build/testentrypoint.sh

ENTRYPOINT [ "/tmp/stripcontrol-app/build/testentrypoint.sh" ]


FROM alpine:3.19 AS builder
WORKDIR /app
COPY ./out/ ./
RUN apkArch="$(apk --print-arch)"; \
    echo "arch: $apkArch"; \
    mv "stripcontrol-app-$apkArch" "stripcontrol-app"

# Start fresh from a smaller image
FROM alpine:3.19 AS finalimage
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/stripcontrol-app ./
COPY configs/config.yml ./
# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/stripcontrol-app", "-c","config.yml"]
