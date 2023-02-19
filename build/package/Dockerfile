FROM golang:1.19-alpine AS build_base

RUN apk add --no-cache git gcc libc-dev sqlite

# Set the Current Working Directory inside the container
WORKDIR /tmp/stripcontrol-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download && go mod tidy

COPY . .

FROM build_base as test
# switch to head instead of latest release due to https://github.com/jstemmer/go-junit-report/issues/138
RUN go install github.com/jstemmer/go-junit-report/v2@bfac3ec01f0cf6ed0dbff9f5f2cb48eb04db38a7 \
   && go install github.com/t-yuki/gocover-cobertura@latest \
# Unit tests
  && chmod +x build/testentrypoint.sh

CMD [ "/tmp/stripcontrol-app/build/testentrypoint.sh" ]


FROM build_base as build
# Build the Go app
RUN go build -o ./out/stripcontrol-app cmd/service/main.go

# Start fresh from a smaller image
FROM alpine:3.17 as finalimage
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=build /tmp/stripcontrol-app/out/stripcontrol-app ./
COPY configs/config.yml ./
# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/stripcontrol-app", "-c","config.yml"]
