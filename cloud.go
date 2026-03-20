package slackhappy

import (
    "bytes"
    "context"
    "contrib.go.opencensus.io/exporter/stackdriver"
    "encoding/json"
    "fmt"
    "github.com/CedricFinance/slackhappy/internal"
    "github.com/CedricFinance/slackhappy/workatoday"
    "github.com/slack-go/slack"
    "go.opencensus.io/plugin/ochttp"
    "go.opencensus.io/trace"
    "log"
    "net/http"
    "os"
    "time"
)

type Request struct {
    Birthdays     bool
    Anniversaries bool
    DryRun        bool
}

var anniversariesWisher *internal.Wisher
var birthdaysWisher *internal.Wisher
var employeesService *workatoday.EmployeeService
var slackNotifier *internal.SlackNotifier

var exporter *stackdriver.Exporter

func MustEnv(name string) string {
    value := os.Getenv(name)

    if value == "" {
        log.Panicf("The environment variable %q is not defined", name)
    }

    return value
}

func SecureEnv(secureName, regularName string) string {
    if value := os.Getenv(secureName); value != "" {
        return value
    }
    return MustEnv(regularName)
}

func init() {
    channelId := MustEnv("SLACK_CHANNEL_ID")
    ownerUserId := MustEnv("OWNER_USER_ID")
    slackToken := SecureEnv("SECURE_SLACK_TOKEN", "SLACK_TOKEN")
    workatodayDomain := MustEnv("WORKATODAY_DOMAIN")
    workatodayToken := SecureEnv("SECURE_WORKATODAY_TOKEN", "WORKATODAY_TOKEN")

    slackClient := slack.New(slackToken, slack.OptionHTTPClient(&http.Client{Transport: &ochttp.Transport{}}))

    slackNotifier = &internal.SlackNotifier{
        SlackClient: slackClient,
        ChannelId:   channelId,
        OwnerId:     ownerUserId,
    }

    anniversariesWisher = &internal.Wisher{
        FilterPredicate: func(employee internal.Employee, date time.Time) bool {
            return employee.IsAnniversary(date)
        },
        Formatter:    internal.SeniorityFormatter{Prefix: ":woop: Happy BlaBl’Anniversary to"},
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

    employeesService = workatoday.NewEmployeeService(
        &workatoday.Config{BaseURL: workatodayDomain, ApiKey: workatodayToken},
        &http.Client{Transport: &ochttp.Transport{}},
    )

    initTracing()
}

func initTracing() {
    var err error

    projectID := os.Getenv("GCP_PROJECT")

    if projectID == "" {
        return
    }

    exporter, err = stackdriver.NewExporter(stackdriver.Options{
        ProjectID: projectID,
        OnError: func(err error) {
            fmt.Printf("Exporter error: %q", err)
        },
    })
    if err != nil {
        log.Panic(err)
    }

    trace.RegisterExporter(exporter)
    trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

type PubSubMessage struct {
    Data []byte `json:"data"`
}

func OnPubSubMessage(ctx context.Context, message PubSubMessage) error {
    ctx, span := trace.StartSpan(ctx, "HappyTrigger")
    defer span.End()
    if exporter != nil {
        defer exporter.Flush()
    }

    var request Request
    decoder := json.NewDecoder(bytes.NewReader(message.Data))
    err := decoder.Decode(&request)
    if err != nil {
        return err
    }

    currentDate := time.Now()
    employees, err := employeesService.ListContext(ctx)
    if err != nil {
        slackNotifier.NotifyError(ctx, fmt.Sprintf("Failed to list employees for birthdays/anniversaries: %v", err))
        return err
    }

    if request.Birthdays {
        message, err := birthdaysWisher.Wish(ctx, currentDate, employees, request.DryRun)
        if err != nil {
            log.Print("Failed to wish birthdays")
            return err
        }
        log.Printf("Today's birthdays message is: %q", message)
    }

    if request.Anniversaries {
        message, err := anniversariesWisher.Wish(ctx, currentDate, employees, request.DryRun)
        if err != nil {
            log.Print("Failed to wish anniversaries")
            return err
        }
        log.Printf("Today's anniversairies message is: %q", message)
    }

    return nil
}
