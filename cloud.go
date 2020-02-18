package slackhappy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/CedricFinance/slackhappy/bamboohr"
	"github.com/CedricFinance/slackhappy/internal"
	"github.com/nlopes/slack"
	"log"
	"os"
	"time"
)

type Request struct {
	Birthdays     bool
	Anniversaries bool
}

var anniversariesWisher *internal.Wisher
var birthdaysWisher *internal.Wisher
var employeesRepository internal.EmployeeRepository

func MustEnv(name string) string {
	value := os.Getenv(name)

	if value == "" {
		log.Panicf("The environment variable %q is not defined", name)
	}

	return value
}

func init() {
	channelId := MustEnv("SLACK_CHANNEL_ID")
	slackToken := MustEnv("SLACK_TOKEN")
	bambooDomain := MustEnv("BAMBOOHR_DOMAIN")
	bambooToken := MustEnv("BAMBOOHR_TOKEN")

	fmt.Printf("%+v", os.Environ())

	slackClient := slack.New(slackToken)

	slackNotifier := &internal.SlackNotifier{
		SlackClient: slackClient,
		ChannelId:   channelId,
	}

	anniversariesWisher = &internal.Wisher{
		FilterPredicate: func(employee internal.Employee, date time.Time) bool {
			return employee.IsAnniversary(date)
		},
		Formatter:    internal.SimpleFormatter{Prefix: "Happy BlaBl’Anniversary to", Suffix: ":woop:!"},
		EmptyMessage: "No anniversaries",
		Notifier:     slackNotifier,
	}

	birthdaysWisher = &internal.Wisher{
		FilterPredicate: func(employee internal.Employee, date time.Time) bool {
			return employee.HireDate.Before(date) && employee.IsBirthday(date)
		},
		Formatter:    internal.SimpleFormatter{Prefix: "Happy Birthday to", Suffix: ":birthday:!"},
		EmptyMessage: "No birthdays",
		Notifier:     slackNotifier,
	}

	employeesRepository = &internal.BambooRepository{
		Client: bamboohr.New(
			bambooDomain,
			bambooToken,
		),
	}

}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func OnPubSubMessage(ctx context.Context, message PubSubMessage) error {
	var request Request
	decoder := json.NewDecoder(bytes.NewReader(message.Data))
	err := decoder.Decode(&request)
	if err != nil {
		return err
	}

	currentDate := time.Now()
	employees := employeesRepository.List(ctx)

	if request.Birthdays {
		message, err := birthdaysWisher.Wish(ctx, currentDate, employees)
		if err != nil {
			log.Print("Failed to wish birthdays")
			return err
		}
		log.Printf("Today's birthdays message is: %q", message)
	}

	if request.Anniversaries {
		message, err := anniversariesWisher.Wish(ctx, currentDate, employees)
		if err != nil {
			log.Print("Failed to wish anniversaries")
			return err
		}
		log.Printf("Today's anniversairies message is: %q", message)
	}

	return nil
}
