> **⚠️ Note:** To use this GitHub Action, you must have access to GitHub Actions. GitHub Actions are currently only available in public beta (you must apply for access).

This action is a part of [GitHub Actions Library](https://github.com/rtCamp/github-actions-library/) created by [rtCamp](https://github.com/rtCamp/).

# Slack Notify - GitHub Action

A [GitHub Action](https://github.com/features/actions) to send a message to a Slack channel.

<img width="485" alt="action-slack-notify-rtcamp" src="https://user-images.githubusercontent.com/4115/54996943-9d38c700-4ff0-11e9-9d35-7e2c16ef0d62.png">

## Usage

You can use this action after any other action. Here is an example setup of this action:

1. Create a `.github/main.workflow` in your GitHub repo, if one doesn't exist already.
2. Add the following code to the `main.workflow` file.

```bash
workflow "Slack Notification Demo" {
  resolves = ["Slack Notification"]
  on = "push"
}

action "Slack Notification" {
  uses = "rtCamp/action-slack-notify@master"
  secrets = ["SLACK_WEBHOOK"]
}
```

3. Create `SLACK_WEBHOOK` secret using [GitHub Action's Secret](https://developer.github.com/actions/creating-workflows/storing-secrets). You can [generate a Slack incoming webhook token from here](https://slack.com/apps/A0F7XDUAZ-incoming-webhooks).


## Environment Variables

By default, action is designed to run with minimal configuration but you can alter Slack notification using following environment variables:

Variable       | Default                                               | Purpose
---------------|-------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------
SLACK_CHANNEL  | Set during Slack webhook creation                     | Specify Slack channel in which message needs to be sent
SLACK_USERNAME | `rtBot`                                               | The name of the sender of the message. Does not need to be a "real" username
SLACK_ICON     | ![rtBot Avatar](https://github.com/rtBot.png?size=32) | User/Bot icon shown with Slack message
SLACK_COLOR    | `good` (green)                                        | You can pass an RGB value like `#efefef` which would change color on left side vertical line of Slack message.
SLACK_MESSAGE  | <Git commit message>                                  | The main Slack message in attachment. It is advised not to override this as this is dynamically generated from Github commit metadata.
SLACK_TITLE    | Message                                               | Title to use before main Slack message

You can see actions with all variables as below:

```bash
workflow "Slack Notification Demo" {
  resolves = ["Slack Notification"]
  on = "push"
}

action "Slack Notification" {
  uses = "rtCamp/action-slack-notify@master"
  env = {
    SLACK_CHANNEL = "general",
    SLACK_USERNAME = "rtCamp",
    SLACK_ICON  = "https://github.com/rtCamp.png?size=48",
    SLACK_COLOR = "#3278BD"
    SLACK_TITLE = "Post Title",
    SLACK_MESSAGE = "Post Content :rocket:"
  }
  secrets = ["SLACK_WEBHOOK"]
}
```

## Hashicorp Vault (Optional)

This GitHub action supports [Hashicorp Vault](https://www.vaultproject.io/). 

To enable Hashicorp Vault support, please define following GitHub secrets:

Variable      | Purpose                                                                       | Example Vaule
--------------|-------------------------------------------------------------------------------|-------------
`VAULT_ADDR`  | [Vault server address](https://www.vaultproject.io/docs/commands/#vault_addr) | `https://example.com:8200`
`VAULT_TOKEN` | [Vault token](https://www.vaultproject.io/docs/concepts/tokens.html)          | `s.gIX5MKov9TUp7iiIqhrP1HgN`

You will need to change `secrets` line in `main.workflow` file to look like below.

```bash
workflow "Slack Notification Demo" {
  resolves = ["Slack Notification"]
  on = "push"
}

action "Slack Notification" {
  uses = "rtCamp/action-slack-notify@master"
  secrets = ["VAULT_ADDR", "VAULT_TOKEN"]
}
```

GitHub action uses `VAULT_TOKEN` to connect to `VAULT_ADDR` to retrieve slack webhook from Vault.

In the Vault, the Slack webhook should be setup as field `webhook` on path `secret/slack`.

## License

[MIT](LICENSE) © 2019 rtCamp
