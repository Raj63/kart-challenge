package slack

import (
	"context"

	"github.com/slack-go/slack"
)

// Messenger defines the contract for sending messages to Slack.
// It allows sending plain text messages with optional attachments.
//
// Implementations of this interface can use different methods of
// communicating with Slack, such as incoming webhooks or the Slack API.
type Messenger interface {
	// SendMessage sends a message to a Slack channel with optional attachments.
	// Returns the message timestamp or ID upon success, or an error if sending fails.
	SendMessage(ctx context.Context, message string, attachments []slack.Attachment) (string, error)
}
