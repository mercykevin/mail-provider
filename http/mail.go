package http

import (
	"net/http"
	"strings"

	"fmt"
	"github.com/nlopes/slack"
	"github.com/open-falcon/mail-provider/config"
	"github.com/toolkits/smtp"
	"github.com/toolkits/web/param"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}
		content := param.MustString(r, "content")
		//send to slack
		api := slack.New(cfg.Slack.Token)
		params := slack.PostMessageParameters{Username: cfg.Slack.Username, AsUser: true}
		attachment := slack.Attachment{
			Pretext: "Alarm that pop from open falcon",
			Color:   "#e11818",
			Text:    content,
			Title:   "Alarm something",
			Fields:  []slack.AttachmentField{slack.AttachmentField{Title: "Priority", Value: "High", Short: false}},
		}
		params.Attachments = []slack.Attachment{attachment}
		channelID, timestamp, slackErr := api.PostMessage(cfg.Slack.Channel, "Error Alert", params)
		if slackErr != nil {
			fmt.Printf("%s\n", slackErr)
		}
		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
		//send to mail
		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		tos = strings.Replace(tos, ",", ";", -1)
		s := smtp.New(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password)
		err := s.SendMail(cfg.Smtp.From, tos, subject, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
