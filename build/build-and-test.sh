#!/bin/bash
set -x
[ -z "$1" ] && echo "service name not given, exiting" && exit 1

SFX=$RANDOM
TEST_IMG_NAME="test-${1}:${SFX}"
TEST_CONTAINER_NAME="test_${1}_${SFX}"
clean_up_container() {
  docker stop -t 20 "$TEST_CONTAINER_NAME" || true 
  docker rm -v -f "$TEST_CONTAINER_NAME" || true
}

clean_up_image() {
  docker rmi -f "$TEST_IMG_NAME" || true
}

clean_up() {
  clean_up_container
  clean_up_image
}

catch_error() {
  clean_up && echo "error during build and test" && exit 1
}

export_file() {
  docker cp "$TEST_CONTAINER_NAME:/$1" "$2"
}
export_test_results() {
  echo "exporting test results"
  export_file "/tmp/stripcontrol-app/report.xml" "../report.xml"
  export_file "/tmp/stripcontrol-app/coverage.xml" "../coverage.xml"
}

export_builds() {
  echo "exporting builds"
  rm -rf ../out || true
  mkdir ../out || true
  export_file "/tmp/stripcontrol-app/out/." "../out/"
}

build_container() {
  echo "building container"
  docker build --no-cache --rm \
    -f package/Dockerfile \
    --target test \
    -t "$TEST_IMG_NAME" ../ || catch_error
}

run_container() {
  echo "running container with target $1"
  docker run --name "$TEST_CONTAINER_NAME" "$TEST_IMG_NAME" "$1"
}

build_container 

run_container "runtests" || catch_error
export_test_results || catch_error
clean_up_container

run_container "runbuild" || catch_error
export_builds || catch_error

clean_up
