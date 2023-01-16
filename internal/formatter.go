package internal

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Formatter interface {
	Format(employees []Employee) string
}

type SimpleFormatter struct {
	Prefix, Suffix string
}

func (f SimpleFormatter) Format(employees []Employee) string {
	names := getNames(employees)
	joinedNames := strings.Join(names, ", ")
	return fmt.Sprintf("%s %s %s", f.Prefix, joinedNames, f.Suffix)
}

func getNames(employees []Employee) []string {
	var names []string
	for _, employee := range employees {
		names = append(names, fmt.Sprintf("%s %s", employee.FirstName, employee.LastName))
	}
	sort.Strings(names)
	return names
}

type SeniorityFormatter struct {
	Prefix string
}

func (f SeniorityFormatter) Format(employees []Employee) string {
	d := groupBySeniority(employees)

	var messages []string
	for _, seniorityAnniversaries := range d {
		names := getNames(seniorityAnniversaries.employees)
		joinedNames := strings.Join(names, ", ")
		messages = append(messages, fmt.Sprintf("%s (%s)", joinedNames, formatSeniority(seniorityAnniversaries.seniority)))
	}

	return fmt.Sprintf("%s %s", f.Prefix, strings.Join(messages, "; "))
}

func groupBySeniority(employees []Employee) []*data {
	var result []*data

	sort.Sort(bySeniority{employees: employees, currentTime: time.Now().UTC()})

	now := time.Now().UTC()

	var currentGroup *data
	for _, employee := range employees {
		if currentGroup == nil || employee.Seniority(now) != currentGroup.seniority {
			currentGroup = &data{
				seniority: employee.Seniority(now),
			}
			result = append(result, currentGroup)
		}

		currentGroup.employees = append(currentGroup.employees, employee)
	}

	return result
}

type bySeniority struct {
	employees   []Employee
	currentTime time.Time
}

func (b bySeniority) Len() int {
	return len(b.employees)
}

func (b bySeniority) Less(i, j int) bool {
	return b.employees[i].Seniority(b.currentTime) < b.employees[j].Seniority(b.currentTime)
}

func (b bySeniority) Swap(i, j int) {
	b.employees[i], b.employees[j] = b.employees[j], b.employees[i]
}

func formatSeniority(seniority int) string {
	if seniority == 1 {
		return "1 year"
	}

	return fmt.Sprintf("%d years", seniority)
}

type data struct {
	seniority int
	employees []Employee
}
