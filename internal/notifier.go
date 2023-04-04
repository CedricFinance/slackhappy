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
    OwnerId     string
}

func (s *SlackNotifier) Notify(ctx context.Context, message string) error {
    _, _, _, err := s.SlackClient.SendMessageContext(ctx, s.ChannelId, slack.MsgOptionText(message, false))
    return err
}

func (s *SlackNotifier) NotifyError(ctx context.Context, message string) error {
    channel, _, _, _ := s.SlackClient.OpenConversationContext(ctx, &slack.OpenConversationParameters{Users: []string{s.OwnerId}})
    _, _, _, err := s.SlackClient.SendMessageContext(ctx, channel.ID, slack.MsgOptionText(message, false))
    return err
}
