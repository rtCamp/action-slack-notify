#!/usr/bin/env bash

# Check required env variables
flag=0
mode="WEBHOOK"
if [[ -z "$SLACK_WEBHOOK" ]]; then
    flag=1
    missing_secret="SLACK_WEBHOOK"
    if [[ -n "$VAULT_ADDR" ]] && [[ -n "$VAULT_TOKEN" ]]; then
        flag=0
    fi
    if [[ -n "$VAULT_ADDR" ]] || [[ -n "$VAULT_TOKEN" ]]; then
        missing_secret="VAULT_ADDR and/or VAULT_TOKEN"
    fi
fi

if [[ "$flag" -eq 1 ]] && [[ -n "$SLACK_TOKEN" || -n "$SLACK_CHANNEL" ]] ; then
    # Basically, if both SLACK_TOKEN and SLACK_CHANNEL are provided, then it's a token mode
    flag=0
    mode="TOKEN"
fi

if [[ "$flag" -eq 1 ]]; then
    echo -e "[\e[0;31mERROR\e[0m] Secret \`$missing_secret\` is missing. Alternatively, a pair of \`SLACK_TOKEN\` and \`SLACK_CHANNEL\` can be provided. Please add it to this action for proper execution.\nRefer https://github.com/rtCamp/action-slack-notify for more information.\n"
    exit 1
fi

export MSG_MODE="$mode"

# custom path for files to override default files
custom_path="$GITHUB_WORKSPACE/.github/slack"
main_script="/main.sh"

if [[ -d "$custom_path" ]]; then
    rsync -av "$custom_path/" /
    chmod +x /*.sh
fi

bash "$main_script"
