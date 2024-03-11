package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	EnvSlackWebhook   = "SLACK_WEBHOOK"
	EnvSlackCustom    = "SLACK_CUSTOM_PAYLOAD"
	EnvSlackIcon      = "SLACK_ICON"
	EnvSlackIconEmoji = "SLACK_ICON_EMOJI"
	EnvSlackChannel   = "SLACK_CHANNEL"
	EnvSlackTitle     = "SLACK_TITLE"
	EnvSlackMessage   = "SLACK_MESSAGE"
	EnvSlackOnSuccess = "SLACK_MESSAGE_ON_SUCCESS"
	EnvSlackOnFailure = "SLACK_MESSAGE_ON_FAILURE"
	EnvSlackOnCancel  = "SLACK_MESSAGE_ON_CANCEL"
	EnvSlackColor     = "SLACK_COLOR"
	EnvSlackUserName  = "SLACK_USERNAME"
	EnvSlackFooter    = "SLACK_FOOTER"
	EnvGithubActor    = "GITHUB_ACTOR"
	EnvGithubRun      = "GITHUB_RUN"
	EnvSiteName       = "SITE_NAME"
	EnvHostName       = "HOST_NAME"
	EnvMinimal        = "MSG_MINIMAL"
	EnvSlackLinkNames = "SLACK_LINK_NAMES"
	EnvThreadTs       = "SLACK_THREAD_TS"
)

type Webhook struct {
	Text        string       `json:"text,omitempty"`
	UserName    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links"`
	Attachments []Attachment `json:"attachments,omitempty"`
	ThreadTs    string       `json:"thread_ts,omitempty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Pretext    string  `json:"pretext,omitempty"`
	Color      string  `json:"color,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func main() {
	endpoint := os.Getenv(EnvSlackWebhook)
	custom_payload := envOr(EnvSlackCustom, "")
	if custom_payload != "" {
		if err := send_raw(endpoint, []byte(custom_payload)); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
			os.Exit(2)
		}
	} else {
		if endpoint == "" {
			fmt.Fprintln(os.Stderr, "URL is required")
			os.Exit(2)
		}
		text := os.Getenv(EnvSlackMessage)
		if text == "" {
			fmt.Fprintln(os.Stderr, "Message is required")
			os.Exit(3)
		}
		if strings.HasPrefix(os.Getenv("GITHUB_WORKFLOW"), ".github") {
			err := os.Setenv("GITHUB_WORKFLOW", "Link to action run.yaml")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to update the workflow's variables: %s\n\n", err)
				os.Exit(4)
			}
		}

		long_sha := os.Getenv("GITHUB_SHA")
		commit_sha := long_sha[0:6]

		color := ""
		switch os.Getenv(EnvSlackColor) {
		case "success":
			color = "good"
			text = envOr(EnvSlackOnSuccess, text) // If exists, override with on success
		case "cancelled":
			color = "#808080"
			text = envOr(EnvSlackOnCancel, text) // If exists, override with on cancelled
		case "failure":
			color = "danger"
			text = envOr(EnvSlackOnFailure, text) // If exists, override with on failure
		default:
			color = envOr(EnvSlackColor, "good")
		}

		if text == "" {
			text = "EOM"
		}

		minimal := os.Getenv(EnvMinimal)
		fields := []Field{}
		if minimal == "true" {
			mainFields := []Field{
				{
					Title: os.Getenv(EnvSlackTitle),
					Value: text,
					Short: false,
				},
			}
			fields = append(mainFields, fields...)
		} else if minimal != "" {
			requiredFields := strings.Split(minimal, ",")
			mainFields := []Field{
				{
					Title: os.Getenv(EnvSlackTitle),
					Value: text,
					Short: false,
				},
			}
			for _, requiredField := range requiredFields {
				switch strings.ToLower(requiredField) {
				case "ref":
					field := []Field{
						{
							Title: "Ref",
							Value: os.Getenv("GITHUB_REF"),
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				case "event":
					field := []Field{
						{
							Title: "Event",
							Value: os.Getenv("GITHUB_EVENT_NAME"),
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				case "actions url":
					field := []Field{
						{
							Title: "Actions URL",
							Value: "<" + os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + "/commit/" + os.Getenv("GITHUB_SHA") + "/checks|" + os.Getenv("GITHUB_WORKFLOW") + ">",
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				case "commit":
					field := []Field{
						{
							Title: "Commit",
							Value: "<" + os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + "/commit/" + os.Getenv("GITHUB_SHA") + "|" + commit_sha + ">",
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				}
			}
			fields = append(mainFields, fields...)
		} else {
			mainFields := []Field{
				{
					Title: "Ref",
					Value: os.Getenv("GITHUB_REF"),
					Short: true,
				}, {
					Title: "Event",
					Value: os.Getenv("GITHUB_EVENT_NAME"),
					Short: true,
				},
				{
					Title: "Actions URL",
					Value: "<" + os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + "/commit/" + os.Getenv("GITHUB_SHA") + "/checks|" + os.Getenv("GITHUB_WORKFLOW") + ">",
					Short: true,
				},
				{
					Title: "Commit",
					Value: "<" + os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + "/commit/" + os.Getenv("GITHUB_SHA") + "|" + commit_sha + ">",
					Short: true,
				},
				{
					Title: os.Getenv(EnvSlackTitle),
					Value: text,
					Short: false,
				},
			}
			fields = append(mainFields, fields...)
		}

		hostName := os.Getenv(EnvHostName)
		if hostName != "" {
			newfields := []Field{
				{
					Title: os.Getenv("SITE_TITLE"),
					Value: os.Getenv(EnvSiteName),
					Short: true,
				},
				{
					Title: os.Getenv("HOST_TITLE"),
					Value: os.Getenv(EnvHostName),
					Short: true,
				},
			}
			fields = append(newfields, fields...)
		}

		msg := Webhook{
			UserName:  os.Getenv(EnvSlackUserName),
			IconURL:   os.Getenv(EnvSlackIcon),
			IconEmoji: os.Getenv(EnvSlackIconEmoji),
			Channel:   os.Getenv(EnvSlackChannel),
			LinkNames: os.Getenv(EnvSlackLinkNames),
			ThreadTs:  os.Getenv(EnvThreadTs),
			Attachments: []Attachment{
				{
					Fallback:   envOr(EnvSlackMessage, "GITHUB_ACTION="+os.Getenv("GITHUB_ACTION")+" \n GITHUB_ACTOR="+os.Getenv("GITHUB_ACTOR")+" \n GITHUB_EVENT_NAME="+os.Getenv("GITHUB_EVENT_NAME")+" \n GITHUB_REF="+os.Getenv("GITHUB_REF")+" \n GITHUB_REPOSITORY="+os.Getenv("GITHUB_REPOSITORY")+" \n GITHUB_WORKFLOW="+os.Getenv("GITHUB_WORKFLOW")),
					Color:      color,
					AuthorName: envOr(EnvGithubActor, ""),
					AuthorLink: os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv(EnvGithubActor),
					AuthorIcon: os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv(EnvGithubActor) + ".png?size=32",
					Footer:     envOr(EnvSlackFooter, "<https://github.com/rtCamp/github-actions-library|Powered By rtCamp's GitHub Actions Library> | <"+os.Getenv(EnvGithubRun)+"|Triggered on this workflow run>"),
					Fields:     fields,
				},
			},
		}

		if err := send(endpoint, msg); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
			os.Exit(1)
		}
	}
	fmt.Fprintf(os.Stdout, "Successfully sent the message!")
}

func envOr(name, def string) string {
	if d, ok := os.LookupEnv(name); ok {
		return d
	}
	return def
}

func send(endpoint string, msg Webhook) error {
	enc, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return send_raw(endpoint, enc)
}

func send_raw(endpoint string, payload []byte) error {
	b := bytes.NewBuffer(payload)
	res, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("Error on message: %s\n", res.Status)
	}
	fmt.Println(res.Status)
	return nil
}
