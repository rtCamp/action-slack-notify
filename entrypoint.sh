#!/usr/bin/env bash

# Check required env variables
flag=0
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

if [[ "$flag" -eq 1 ]]; then
    printf "[\e[0;31mERROR\e[0m] Secret \`$missing_secret\` is missing. Please add it to this action for proper execution.\nRefer https://github.com/rtCamp/action-slack-notify for more information.\n"
    exit 1
fi

export GITHUB_BRANCH=${GITHUB_REF##*heads/}
export SLACK_ICON=${SLACK_ICON:-"https://avatars0.githubusercontent.com/u/43742164"}
export SLACK_USERNAME=${SLACK_USERNAME:-"rtBot"}
export CI_SCRIPT_OPTIONS="ci_script_options"
export SLACK_TITLE=${SLACK_TITLE:-"Message"}
export COMMIT_MESSAGE=$(cat "/github/workflow/event.json" | jq .commits | jq '.[0].message' -r)

hosts_file="$GITHUB_WORKSPACE/.github/hosts.yml"

if [[ -z "$SLACK_CHANNEL" ]]; then
	user_slack_channel=$(cat "$hosts_file" | shyaml get-value "$CI_SCRIPT_OPTIONS.slack-channel" | tr '[:upper:]' '[:lower:]')
fi

if [[ -n "$user_slack_channel" ]]; then
	export SLACK_CHANNEL="$user_slack_channel"
fi

# Login to vault using GH Token
if [[ -n "$VAULT_GITHUB_TOKEN" ]]; then
	unset VAULT_TOKEN
	vault login -method=github token="$VAULT_GITHUB_TOKEN" > /dev/null
fi

if [[ -n "$VAULT_GITHUB_TOKEN" ]] || [[ -n "$VAULT_TOKEN" ]]; then
	export SLACK_WEBHOOK=$(vault read -field=webhook secret/slack)
fi

if [[ -f "$hosts_file" ]]; then
	hostname=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.hostname")
	user=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.user")
	export HOST_NAME="\`$user@$hostname\`"
	export DEPLOY_PATH=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.deploy_path")

	temp_url=${DEPLOY_PATH%%/app*}
	export SITE_NAME="${temp_url##*sites/}"

	if [[ -n "$SITE_NAME" ]]; then
		export SITE_TITLE="Site"
	fi
	if [[ -n "$HOST_NAME" ]]; then
		export HOST_TITLE="SSH Host"
	fi
fi

if [[ -z "$SLACK_MESSAGE" ]]; then
	export SLACK_MESSAGE="$COMMIT_MESSAGE"
fi

slack-notify "$@"
