#!/usr/bin/env bash

docker build -t ghcr.io/radical-egg/steamcmd-healthcheck:latest \
    --target export .
