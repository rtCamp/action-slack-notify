This action is a part of [GitHub Actions Library](https://github.com/rtCamp/github-actions-library/) created by [rtCamp](https://github.com/rtCamp/).

# Slack Notify - GitHub Action

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

A [GitHub Action](https://github.com/features/actions) to send a message to a Slack channel.

**Screenshot**

<img width="485" alt="action-slack-notify-rtcamp" src="https://user-images.githubusercontent.com/4115/54996943-9d38c700-4ff0-11e9-9d35-7e2c16ef0d62.png">

The `Site` and `SSH Host` details are only available if this action is run after [Deploy WordPress GitHub action](https://github.com/rtCamp/action-deploy-wordpress).

## Usage

You can use this action after any other action. Here is an example setup of this action:

1. Create a `.github/workflows/slack-notify.yml` file in your GitHub repo.
2. Add the following code to the `slack-notify.yml` file.

```yml
on: push
name: Slack Notification Demo
jobs:
  slackNotification:
    name: Slack Notification
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Slack Notification
      uses: rtCamp/action-slack-notify@v2
      env:
        SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
```

### 📢 Send Notification to Multiple Slack Channels

You can notify multiple Slack channels by providing a comma-separated list of channel names:

```yaml
jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - uses: rtCamp/action-slack-notify@v2
        env:
          SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
          SLACK_CHANNEL: "#general,#build-status"
          SLACK_MESSAGE: "🚀 Deployment completed successfully!"
```

### 💡 Advanced: Send Different Messages to Different Channels (Matrix Strategy)

If you'd like to send **custom messages to different Slack channels**, use GitHub Actions' matrix strategy. Here's how:

```yaml
on: push
name: Slack Notification Matrix Example

jobs:
  notify:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - channel: '#build'
            message: '✅ Build completed successfully!'
          - channel: '#deployments'
            message: '🚀 Deployment finished without errors!'
    steps:
      - uses: actions/checkout@v4
      - name: Notify Slack
        uses: rtCamp/action-slack-notify@v2
        env:
          SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
          SLACK_CHANNEL: ${{ matrix.channel }}
          SLACK_MESSAGE: ${{ matrix.message }}
```

This setup will send **two different messages** to two different channels. It's ideal when you want CI messages in `#build` and deployment updates in `#deployments`.

3. Create `SLACK_WEBHOOK` secret using [GitHub Action's Secret](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets#creating-encrypted-secrets-for-a-repository). You can [generate a Slack incoming webhook token from here](https://slack.com/apps/A0F7XDUAZ-incoming-webhooks).

## Environment Variables

By default, action is designed to run with minimal configuration but you can alter Slack notification using following environment variables:

| Variable                    | Default                                               | Purpose                                                                                                                                                                                                                                                                                                                  |
| --------------------------- | ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| SLACK\_CHANNEL              | Set during Slack webhook creation                     | Specify Slack channel in which message needs to be sent                                                                                                                                                                                                                                                                  |
| SLACK\_USERNAME             | `rtBot`                                               | Custom Slack Username sending the message. Does not need to be a "real" username.                                                                                                                                                                                                                                        |
| SLACK\_MSG\_AUTHOR          | `$GITHUB_ACTOR` (The person who triggered action).    | GitHub username of the person who has triggered the action. In case you want to modify it, please specify correct GitHub username.                                                                                                                                                                                       |
| SLACK\_ICON                 | ![rtBot Avatar](https://github.com/rtBot.png?size=32) | User/Bot icon shown with Slack message. It uses the URL supplied to this env variable to display the icon in slack message.                                                                                                                                                                                              |
| SLACK\_ICON\_EMOJI          | -                                                     | User/Bot icon shown with Slack message, in case you do not wish to add a URL for slack icon as above, you can set slack emoji in this env variable. Example value: `:bell:` or any other valid slack emoji.                                                                                                              |
| SLACK\_COLOR                | `good` (green)                                        | You can pass `${{ job.status }}` for automatic coloring or an RGB value like `#efefef` which would change color on left side vertical line of Slack message. Other valid values for this field are: `success`, `cancelled` or `failure`.                                                                                 |
| SLACK\_LINK\_NAMES          | -                                                     | If set to `true`, enable mention in Slack message.                                                                                                                                                                                                                                                                       |
| SLACK\_MESSAGE              | Generated from git commit message.                    | The main Slack message in attachment. It is advised not to override this.                                                                                                                                                                                                                                                |
| SLACK\_TITLE                | Message                                               | Title to use before main Slack message.                                                                                                                                                                                                                                                                                  |
| SLACK\_FOOTER               | Powered By rtCamp's GitHub Actions Library            | Slack message footer.                                                                                                                                                                                                                                                                                                    |
| MSG\_MINIMAL                | -                                                     | If set to `true`, removes: `Ref`, `Event`,  `Actions URL` and `Commit` from the message. You can optionally whitelist any of these 4 removed values by passing it comma separated to the variable instead of `true`. (ex: `MSG_MINIMAL: event` or `MSG_MINIMAL: ref,actions url`, etc.)                                  |
| SLACKIFY\_MARKDOWN          | -                                                     | If set to `true`, it will convert markdown to slack format. (ex: `*bold*` to `bold`) Note: This only works for custom messages and not for the default message generated by the action. Credits: [slackify-markdown-action](https://github.com/marketplace/actions/slack-markdown-converter)                             |
| SLACK\_THREAD\_TS           | -                                                     | If you want to send message in a thread, you can pass the timestamp of the parent message to this variable. You can get the timestamp of the parent message from the message URL in Slack. (ex: `SLACK_THREAD_TS: 1586130833.000100`)                                                                                    |
| SLACK\_TOKEN                | -                                                     | If you want to send message to a channel using a slack token. You will need to pass a channel in order to send messages using token, requiring a value for `SLACK_CHANNEL`. Note that in case both webhook url and token are provided, webhook url will be prioritized.                                                  |
| SLACK\_MESSAGE\_ON\_SUCCESS | -                                                     | If set, will send the provided message instead of the default message when the passed status (through `SLACK_COLOR`) is `success`.                                                                                                                                                                                       |
| SLACK\_MESSAGE\_ON\_FAILURE | -                                                     | If set, will send the provided message instead of the default message when the passed status (through `SLACK_COLOR`) is `failure`.                                                                                                                                                                                       |
| SLACK\_MESSAGE\_ON\_CANCEL  | -                                                     | If set, will send the provided message instead of the default message when the passed status (through `SLACK_COLOR`) is `cancelled`.                                                                                                                                                                                     |
| SLACK\_CUSTOM\_PAYLOAD      | -                                                     | If you want to send a custom payload to slack, you can pass it as a string to this variable. This will override all other variables and send the custom payload to slack. Example: `SLACK_CUSTOM_PAYLOAD: '{"text": "Hello, World!"}'`, Note: This payload should be in JSON format, and is not validated by the action. |
| SLACK\_FILE\_UPLOAD         | -                                                     | If you want to upload a file to slack, you can pass the file path to this variable. Example: `SLACK_FILE_UPLOAD: /path/to/file.txt`. Note: This file should be present in the repository, or github workspace. Otherwise, should be accessable in the container the action is running in.                                |
| ENABLE\_ESCAPES             | -                                                     | If set to `true`, will enable backslash escape sequences such as `\n`, `\t`, etc. in the message. Note: This only works for custom messages and not for the default message generated by the action.                                                                                                                     |

You can see the action block with all variables as below:

```yml
    - name: Slack Notification
      uses: rtCamp/action-slack-notify@v2
      env:
        SLACK_CHANNEL: general
        SLACK_COLOR: ${{ job.status }} # or a specific color like 'good' or '#ff00ff'
        SLACK_ICON: https://github.com/rtCamp.png?size=48
        SLACK_MESSAGE: 'Post Content :rocket:'
        SLACK_TITLE: Post Title
        SLACK_USERNAME: rtCamp
        SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
```

Below screenshot help you visualize message part controlled by different variables:

<img width="600" alt="Screenshot_2019-03-26_at_5_56_05_PM" src="https://user-images.githubusercontent.com/4115/54997488-d1f94e00-4ff1-11e9-897f-a35ab90f525f.png">

The `Site` and `SSH Host` details are only available if this action is run after [Deploy WordPress GitHub action](https://github.com/rtCamp/action-deploy-wordpress).

## Hashicorp Vault (Optional) (Deprecated)

This GitHub action supports [Hashicorp Vault](https://www.vaultproject.io/).

To enable Hashicorp Vault support, please define following GitHub secrets:

| Variable      | Purpose                                                                       | Example Vaule                |
| ------------- | ----------------------------------------------------------------------------- | ---------------------------- |
| `VAULT_ADDR`  | [Vault server address](https://www.vaultproject.io/docs/commands/#vault_addr) | `https://example.com:8200`   |
| `VAULT_TOKEN` | [Vault token](https://www.vaultproject.io/docs/concepts/tokens.html)          | `s.gIX5MKov9TUp7iiIqhrP1HgN` |

You will need to change `secrets` line in `slack-notify.yml` file to look like below.

```yml
on: push
name: Slack Notification Demo
jobs:
  slackNotification:
    name: Slack Notification
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Slack Notification
      uses: rtCamp/action-slack-notify@v2
      env:
        VAULT_ADDR: ${{ secrets.VAULT_ADDR }}
        VAULT_TOKEN: ${{ secrets.VAULT_TOKEN }}
```

GitHub action uses `VAULT_TOKEN` to connect to `VAULT_ADDR` to retrieve slack webhook from Vault.

In the Vault, the Slack webhook should be setup as field `webhook` on path `secret/slack`.

## Credits

Source: [technosophos/slack-notify](https://github.com/technosophos/slack-notify)

## License

[MIT](LICENSE) © 2022 rtCamp

## Does this interest you?

<a href="https://rtcamp.com/"><img src="https://rtcamp.com/wp-content/uploads/sites/2/2019/04/github-banner@2x.png" alt="Join us at rtCamp, we specialize in providing high performance enterprise WordPress solutions"></a>
