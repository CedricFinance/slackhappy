package internal

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type Wisher struct {
	FilterPredicate func(employee Employee, date time.Time) bool
	Formatter       Formatter
	EmptyMessage    string
	Notifier        Notifier
}

func (w Wisher) Wish(ctx context.Context, date time.Time, employees []Employee) (string, error) {
	filtered := Filter(employees, func(employee Employee) bool {
		return w.FilterPredicate(employee, date)
	})

	if len(filtered) == 0 {
		return w.EmptyMessage, nil
	}

	names := GetNames(filtered)

	message := w.Formatter.Format(names)

	err := w.Notifier.Notify(ctx, message)
	if err != nil {
		return "", err
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
