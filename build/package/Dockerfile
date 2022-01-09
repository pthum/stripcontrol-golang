FROM golang:1.17-alpine AS build_base

RUN apk add --no-cache git gcc libc-dev sqlite

# Set the Current Working Directory inside the container
WORKDIR /tmp/stripcontrol-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download && go mod tidy

COPY . .

# Unit tests
# RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/stripcontrol-app cmd/service/main.go

# Start fresh from a smaller image
FROM alpine:3.12
WORKDIR /app
RUN apk add ca-certificates

COPY --from=build_base /tmp/stripcontrol-app/out/stripcontrol-app ./
COPY static ./static
COPY configs/config.yml ./
# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/stripcontrol-app", "-c","config.yml"]