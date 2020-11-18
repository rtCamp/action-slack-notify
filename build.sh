#!/usr/bin/env bash

docker build -t clouddrove/slack-notify:1.0 .

if [[ $? != 0 ]]; then
    echo "slack-notify docker Build failed."
    exit 1
fi

docker push clouddrove/slack-notify:1.0
