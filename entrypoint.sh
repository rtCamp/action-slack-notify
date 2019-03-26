#!/usr/bin/env bash

export GITHUB_BRANCH=${GITHUB_REF##*heads/}
export SLACK_ICON=${SLACK_ICON:-"https://avatars0.githubusercontent.com/u/43742164"}
export SLACK_USERNAME=${SLACK_USERNAME:-"rtBot"}
export CI_SCRIPT_OPTIONS="ci_script_options"
export SLACK_TITLE=${SLACK_TITLE:-"Message"}

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
	deploy_path=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.deploy_path")

	temp_url=${deploy_path%%/app*}
	site="${temp_url##*sites/}"

	if [[ -n "$site" ]]; then
		export TEMPLATE_SLACK_MESSAGE="Deployed successfully on site: \`$site\` for branch \`$GITHUB_BRANCH\` :tada: on server: \`$hostname\` :rocket:"
	else
		export TEMPLATE_SLACK_MESSAGE="Deployed successfully on \`$hostname\` for branch \`$GITHUB_BRANCH\` :tada: on path: \`$deploy_path\` :rocket:"
	fi
fi

export SLACK_MESSAGE=${SLACK_MESSAGE:-$TEMPLATE_SLACK_MESSAGE}

if [[ -z "$SLACK_MESSAGE" ]]; then
	export SLACK_MESSAGE='Success :tada:'
fi

slack-notify "$@"
