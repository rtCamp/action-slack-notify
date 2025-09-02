package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	EnvSlackUpload    = "SLACK_FILE_UPLOAD"
	EnvMessageMode    = "MSG_MODE"
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
	endpoint := getEnv(EnvSlackWebhook)
	custom_payload := envOr(EnvSlackCustom, "")
	if endpoint == "" {
		if getEnv(EnvSlackChannel) == "" {
			fmt.Fprintln(os.Stderr, "Channel is required for sending message using a token")
			os.Exit(1)
		}
		if getEnv(EnvMessageMode) == "TOKEN" {
			endpoint = "https://slack.com/api/chat.postMessage"
		} else {
			fmt.Fprintln(os.Stderr, "URL is required")
			os.Exit(2)
		}
	}
	if custom_payload != "" {
		if err := send_raw(endpoint, []byte(custom_payload)); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
			os.Exit(2)
		}
	} else {
		text := getEnv(EnvSlackMessage)
		if text == "" {
			fmt.Fprintln(os.Stderr, "Message is required")
			os.Exit(3)
		}
		if strings.HasPrefix(getEnv("GITHUB_WORKFLOW"), ".github") {
			err := os.Setenv("GITHUB_WORKFLOW", "Link to action run.yaml")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to update the workflow's variables: %s\n\n", err)
				os.Exit(4)
			}
		}

		long_sha := getEnv("GITHUB_SHA")
		commit_sha := long_sha[0:6]

		color := ""
		switch strings.ToLower(getEnv(EnvSlackColor)) {
		case "success":
			color = "good"
			// If exists, override with on success
			success_msg := envOr(EnvSlackOnSuccess, "")
			if success_msg != "" {
				text = success_msg
			}
		case "cancelled":
			color = "#808080"
			// If exists, override with on cancel
			cancel_msg := envOr(EnvSlackOnCancel, "")
			if cancel_msg != "" {
				text = cancel_msg
			}
		case "failure":
			color = "danger"
			// If exists, override with on failure
			failure_msg := envOr(EnvSlackOnFailure, "")
			if failure_msg != "" {
				text = failure_msg
			}
		default:
			color = envOr(EnvSlackColor, "good")
		}

		if text == "" {
			text = "EOM"
		}

		minimal := getEnv(EnvMinimal)
		fields := []Field{}
		if minimal == "true" {
			mainFields := []Field{
				{
					Title: getEnv(EnvSlackTitle),
					Value: text,
					Short: false,
				},
			}
			fields = append(mainFields, fields...)
		} else if minimal != "" {
			requiredFields := strings.Split(minimal, ",")
			mainFields := []Field{
				{
					Title: getEnv(EnvSlackTitle),
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
							Value: getEnv("GITHUB_REF"),
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				case "event":
					field := []Field{
						{
							Title: "Event",
							Value: getEnv("GITHUB_EVENT_NAME"),
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				case "actions url":
					field := []Field{
						{
							Title: "Actions URL",
							Value: "<" + getEnv("GITHUB_SERVER_URL") + "/" + getEnv("GITHUB_REPOSITORY") + "/commit/" + getEnv("GITHUB_SHA") + "/checks|" + getEnv("GITHUB_WORKFLOW") + ">",
							Short: true,
						},
					}
					mainFields = append(field, mainFields...)
				case "commit":
					field := []Field{
						{
							Title: "Commit",
							Value: "<" + getEnv("GITHUB_SERVER_URL") + "/" + getEnv("GITHUB_REPOSITORY") + "/commit/" + getEnv("GITHUB_SHA") + "|" + commit_sha + ">",
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
					Value: getEnv("GITHUB_REF"),
					Short: true,
				}, {
					Title: "Event",
					Value: getEnv("GITHUB_EVENT_NAME"),
					Short: true,
				},
				{
					Title: "Actions URL",
					Value: "<" + getEnv("GITHUB_SERVER_URL") + "/" + getEnv("GITHUB_REPOSITORY") + "/commit/" + getEnv("GITHUB_SHA") + "/checks|" + getEnv("GITHUB_WORKFLOW") + ">",
					Short: true,
				},
				{
					Title: "Commit",
					Value: "<" + getEnv("GITHUB_SERVER_URL") + "/" + getEnv("GITHUB_REPOSITORY") + "/commit/" + getEnv("GITHUB_SHA") + "|" + commit_sha + ">",
					Short: true,
				},
				{
					Title: getEnv(EnvSlackTitle),
					Value: text,
					Short: false,
				},
			}
			fields = append(mainFields, fields...)
		}

		hostName := getEnv(EnvHostName)
		if hostName != "" {
			newfields := []Field{
				{
					Title: getEnv("SITE_TITLE"),
					Value: getEnv(EnvSiteName),
					Short: true,
				},
				{
					Title: getEnv("HOST_TITLE"),
					Value: getEnv(EnvHostName),
					Short: true,
				},
			}
			fields = append(newfields, fields...)
		}

		msg := Webhook{
			UserName:  getEnv(EnvSlackUserName),
			IconURL:   getEnv(EnvSlackIcon),
			IconEmoji: getEnv(EnvSlackIconEmoji),
			Channel:   getEnv(EnvSlackChannel),
			LinkNames: getEnv(EnvSlackLinkNames),
			ThreadTs:  getEnv(EnvThreadTs),
			Attachments: []Attachment{
				{
					Fallback:   envOr(EnvSlackMessage, "GITHUB_ACTION="+getEnv("GITHUB_ACTION")+" \n GITHUB_ACTOR="+getEnv("GITHUB_ACTOR")+" \n GITHUB_EVENT_NAME="+getEnv("GITHUB_EVENT_NAME")+" \n GITHUB_REF="+getEnv("GITHUB_REF")+" \n GITHUB_REPOSITORY="+getEnv("GITHUB_REPOSITORY")+" \n GITHUB_WORKFLOW="+getEnv("GITHUB_WORKFLOW")),
					Color:      color,
					AuthorName: envOr(EnvGithubActor, ""),
					AuthorLink: getEnv("GITHUB_SERVER_URL") + "/" + getEnv(EnvGithubActor),
					AuthorIcon: getEnv("GITHUB_SERVER_URL") + "/" + getEnv(EnvGithubActor) + ".png?size=32",
					// Change link to org controlled one.
					Footer: envOr(EnvSlackFooter, "<https://github.com/rtCamp/github-actions-library|Powered By rtCamp's GitHub Actions Library> | <"+getEnv(EnvGithubRun)+"|Triggered on this workflow run>"),
					Fields: fields,
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

func getEnv(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func envOr(name, def string) string {
	if d, ok := os.LookupEnv(name); ok {
		return strings.TrimSpace(d)
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

	var res *http.Response
	var err error

	switch getEnv(EnvMessageMode) {
	case "WEBHOOK":
		res, err = http.Post(endpoint, "application/json", b)
	case "TOKEN":
		req, err := http.NewRequest("POST", endpoint, b)
		if err != nil {
			return fmt.Errorf("Error creating request: %s\n", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getEnv("SLACK_TOKEN"))
		client := &http.Client{}
		res, err = client.Do(req)
	default:
		fmt.Fprintf(os.Stderr, "Invalid message mode: %s\n", getEnv(EnvMessageMode))
		os.Exit(6)
	}

	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("Error on message: %s\n", res.Status)
	}

	if os.Getenv(EnvSlackUpload) != "" {
		err = sendFile(os.Getenv(EnvSlackUpload), "", os.Getenv(EnvSlackChannel), os.Getenv(EnvThreadTs))
		if err != nil {
			return err
		}
	}

	return nil
}

func sendFile(filename string, message string, channel string, thread_ts string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileData := &bytes.Buffer{}
	writer := multipart.NewWriter(fileData)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)

	err = writer.WriteField("initial_comment", message)
	if err != nil {
		return err
	}

	err = writer.WriteField("channels", channel)
	if err != nil {
		return err
	}

	if thread_ts != "" {
		err = writer.WriteField("thread_ts", thread_ts)
		if err != nil {
			return err
		}
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", fileData)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SLACK_TOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("Error on message: %s\n", res.Status)
	}
	fmt.Println(res.Status)
	return nil
}
