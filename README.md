# Slack Notify - GitHub Action

A [GitHub Action](https://github.com/features/actions) that can be used to send a message to a Slack channel.

![Slack Notify Action demo](https://user-images.githubusercontent.com/8456197/54735458-5439da80-4bcc-11e9-9931-393e32c79ad1.png)

This action is a part of [GitHub action library](https://github.com/rtCamp/github-actions-library/) created by [rtCamp](https://github.com/rtCamp/).

## Installation

To use this GitHub Action, you must have access to GitHub Actions. GitHub Actions are currently only available in public beta (you must [apply for access](https://github.com/features/actions)).

You can use this action after any other action to send success or failure notification on Slack. Here is an example setup of this action:

1. Create a `.github/main.workflow` in your GitHub repo.
2. Add the following code to the `main.workflow` file and commit it to the repo's `master` branch.

```bash
workflow "" {
  resolves = ["Slack Notification"]
  on = "push"
}

action "Slack Notification" {
  uses = "rtCamp/action-slack-notify@master"
  env = {
    SLACK_MESSAGE = "Commit received :rocket:",
    SLACK_USERNAME = "bot-account"
  }
  secrets = ["SLACK_WEBHOOK"]
}
```

3. Define `SLACK_WEBHOOK` as a [GitHub Actions Secret](https://developer.github.com/actions/creating-workflows/storing-secrets). (You can add secrets using the visual workflow editor or the repository settings.)
4. Whenever you commit, this action will run.

## Environment Variables

```bash
# The Slack-assigned webhook
SLACK_WEBHOOK=https://hooks.slack.com/services/Txxxxxx/Bxxxxxx/xxxxxxxx
# A URL to an icon
SLACK_ICON=http://example.com/icon.png
# The channel to send the message to (if omitted, use Slack-configured default)
SLACK_CHANNEL=example
# The title of the message
SLACK_TITLE="Hello World"
# The body of the message
SLACK_MESSAGE="Today is a fine day"
# RGB color to for message formatting. (Slack determines what is colored by this)
SLACK_COLOR="#efefef"
# The name of the sender of the message. Does not need to be a "real" username
SLACK_USERNAME="notify-bot"
```
## Additional Vault Support

Instead of setting up `secrets = ["SLACK_WEBHOOK"]`, slack webhook can also be read from [vault](https://www.vaultproject.io/). In vault, the slack webhook should be setup as field `webhook` on path `secret/slack`.

`VAULT_ADDR` secret variable specifies the address on which vault is deployed, e.g., `VAULT_ADDR=https://example.com:8200`. [VAULT_TOKEN](https://www.vaultproject.io/docs/concepts/tokens.html) is the token by which authentication with vault will be possible to retrieve the secrets and information.

Example usage with vault setup:

```bash
workflow "" {
  resolves = ["Slack Notification"]
  on = "push"
}

action "Slack Notification" {
  uses = "rtCamp/action-slack-notify@master"
  env = {
    SLACK_MESSAGE = "Commit received :rocket:",
    SLACK_USERNAME = "bot-account"
  }
  secrets = ["VAULT_ADDR", "VAULT_TOKEN"]
}
```

## License

[MIT](LICENSE) © 2019 rtCamp
