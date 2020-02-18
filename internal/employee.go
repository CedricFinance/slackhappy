package internal

import "time"

type Employee struct {
	FirstName string
	LastName  string
	Birthday  string
	HireDate  time.Time
}

func (e Employee) IsAnniversary(date time.Time) bool {
	return e.HireDate.Year() != date.Year() &&
		e.HireDate.Month() == date.Month() &&
		e.HireDate.Day() == date.Day()
}

func (e Employee) IsBirthday(date time.Time) bool {
	monthAndDay := date.Format("01-02")

	return e.Birthday == monthAndDay
}
