package slack

import (
	"context"
	"fmt"
	"library/logger"
	"log"

	"github.com/slack-go/slack"
)

type slackMessenger struct {
	botToken  string
	channelID string
	client    *slack.Client
	logger    *logger.Logger
}

// NewSlackMessanger creates a new Messenger instance for sending messages to Slack.
// It requires a bot token, channel ID, and a logger. Returns a Messenger or an error if parameters are invalid.
func NewSlackMessanger(botToken, channelID string, logger *logger.Logger) (Messenger, error) {

	if botToken == "" || channelID == "" {
		return nil, fmt.Errorf("channelId or slack bottoken cannot be empty")
	}
	return &slackMessenger{
		botToken:  botToken,
		channelID: channelID,
		client:    slack.New(botToken),
		logger:    logger,
	}, nil

}

// SendMessage sends a message with optional attachments to the configured Slack channel.
// It returns the timestamp of the sent message or an error if sending fails.
func (s *slackMessenger) SendMessage(ctx context.Context, message string, attachments []slack.Attachment) (string, error) {
	// Send a new message if there's no thread
	channel, timestamp, err := s.client.PostMessage(
		s.channelID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(attachments...),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("New message sent to channel %s with timestamp %s", channel, timestamp)
	return timestamp, nil
}
