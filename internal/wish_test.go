package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var employees = []Employee{
	{FirstName: "John", LastName: "Doe", Birthday: "01-02", HireDate: MustParseTime("2006-01-02", "2006-01-02")},
	{FirstName: "Emily", LastName: "Brow", Birthday: "03-14", HireDate: MustParseTime("2006-01-02", "2007-10-12")},
	{FirstName: "Stephanie", LastName: "Wilson", Birthday: "03-14", HireDate: MustParseTime("2006-01-02", "2008-12-31")},
}

type FakeNotifier struct {
	Message string
}

func (f *FakeNotifier) Notify(ctx context.Context, message string) error {
	f.Message = message
	return nil
}

func TestWisher_Wish_EmptyMessage(t *testing.T) {
	wisher := Wisher{
		FilterPredicate: func(employee Employee, date time.Time) bool {
			return false
		},
		Formatter:    nil,
		EmptyMessage: "Empty",
		Notifier:     nil,
	}

	result, err := wisher.Wish(context.Background(), time.Now(), employees)

	assert.Nil(t, err)
	assert.Equal(t, "Empty", result)
}

func TestWisher_Wish_All(t *testing.T) {
	fakeNotifier := &FakeNotifier{}
	wisher := Wisher{
		FilterPredicate: func(employee Employee, date time.Time) bool {
			return true
		},
		Formatter:    SimpleFormatter{Prefix: "Prefix", Suffix: "Suffix"},
		EmptyMessage: "Empty",
		Notifier:     fakeNotifier,
	}

	result, err := wisher.Wish(context.Background(), time.Now(), employees)

	assert.Nil(t, err)
	assert.Equal(t, "Prefix Emily Brow, John Doe, Stephanie Wilson Suffix", result)
	assert.Equal(t, "Prefix Emily Brow, John Doe, Stephanie Wilson Suffix", fakeNotifier.Message)
}

func TestWisher_Wish_Birthdays(t *testing.T) {
	fakeNotifier := &FakeNotifier{}
	wisher := Wisher{
		FilterPredicate: func(employee Employee, date time.Time) bool {
			return employee.IsBirthday(date)
		},
		Formatter:    SimpleFormatter{Prefix: "Prefix", Suffix: "Suffix"},
		EmptyMessage: "Empty",
		Notifier:     fakeNotifier,
	}

	result, err := wisher.Wish(context.Background(), MustParseTime("2006-01-02", "2020-03-14"), employees)

	assert.Nil(t, err)
	assert.Equal(t, "Prefix Emily Brow, Stephanie Wilson Suffix", result)
	assert.Equal(t, "Prefix Emily Brow, Stephanie Wilson Suffix", fakeNotifier.Message)
}
