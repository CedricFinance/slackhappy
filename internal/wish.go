package internal

import (
	"context"
	"fmt"
	"log"
	"sort"
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

	names := GetNames(filtered)

	message := w.Formatter.Format(names)

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

func GetNames(employees []Employee) []string {
	var names []string
	for _, employee := range employees {
		names = append(names, fmt.Sprintf("%s %s", employee.FirstName, employee.LastName))
	}
	sort.Strings(names)
	return names
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
