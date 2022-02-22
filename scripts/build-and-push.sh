#!/bin/bash

IMAGE=$1
BUILDX_INSTANCE="action-slack-notify-builder"
PLATFORMS="linux/arm64,linux/amd64"

function print_usage {
  echo "$0 <Image>"
}

function validate_command_line_args {
  if [[ -z "$IMAGE" ]]; then 
    echo "Image is not supplied"
    print_usage
    exit 1
  fi

  echo "=================================="
  echo "Arguments"
  echo "Image: $IMAGE"
  echo "=================================="
}

function validate_docker_installed {
  if ! [[ $(command -v docker) ]]; then 
    echo "Docker is not installed. Install docker first"
    exit 1
  fi
}

function setup_buildx_platforms {
  echo "Install buildx platforms"
  docker run --privileged --rm tonistiigi/binfmt --install all > /dev/null
}

function setup_buildx_instance {
  echo "Check buildx instance '$BUILDX_INSTANCE' is set"

  if [[ $(docker buildx ls | grep $BUILDX_INSTANCE | grep docker-container | wc -l ) -eq 0 ]]; then
    echo "Create buildx instance '$BUILDX_INSTANCE'"
    docker buildx create --name=$BUILDX_INSTANCE --driver=docker-container
    echo "Successfully created and converted to buildx instance '$BUILDX_INSTANCE'"
  else
    echo "Buildx instance $BUILDX_INSTANCE is ready. Skip..."
  fi
}

function build_and_push_images {
  echo "Build and push docker images"
  docker buildx build --builder=$BUILDX_INSTANCE --platform=$PLATFORMS --push -t $IMAGE .
}


validate_command_line_args
validate_docker_installed

setup_buildx_platforms
setup_buildx_instance

build_and_push_images