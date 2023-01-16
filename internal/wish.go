package internal

import (
	"context"
	"log"
	"time"
)

type Wisher struct {
	FilterPredicate func(employee Employee, date time.Time) bool
	Formatter       Formatter
	EmptyMessage    string
	Notifier        Notifier
}

func (w Wisher) Wish(ctx context.Context, date time.Time, employees []Employee, dryRun bool) (string, error) {
	filtered := Filter(employees, func(employee Employee) bool {
		return w.FilterPredicate(employee, date)
	})

	if len(filtered) == 0 {
		return w.EmptyMessage, nil
	}

	message := w.Formatter.Format(filtered)

	if dryRun {
		log.Println("Dry-run: no message sent")
	} else {
		log.Printf("Sending message: %q", message)
		err := w.Notifier.Notify(ctx, message)
		if err != nil {
			return "", err
		}
	}

	return message, nil
}

func Filter(employees []Employee, predicate func(employee Employee) bool) []Employee {
	var results []Employee

	for _, employee := range employees {
		if predicate(employee) {
			results = append(results, employee)
		}
	}

	return results
}
