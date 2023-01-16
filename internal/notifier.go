package internal

import (
    "context"
    "github.com/slack-go/slack"
)

type Notifier interface {
    Notify(ctx context.Context, message string) error
}

type SlackNotifier struct {
    SlackClient *slack.Client
    ChannelId   string
}

func (s *SlackNotifier) Notify(ctx context.Context, message string) error {
    _, _, _, err := s.SlackClient.SendMessageContext(ctx, s.ChannelId, slack.MsgOptionText(message, false))
    return err
}
