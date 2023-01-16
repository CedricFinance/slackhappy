package internal

import "time"

type Employee struct {
    FirstName string
    LastName  string
    Birthday  time.Time
    HireDate  time.Time
}

func (e Employee) IsAnniversary(date time.Time) bool {
    return e.HireDate.Year() != date.Year() &&
        e.HireDate.Month() == date.Month() &&
        e.HireDate.Day() == date.Day()
}

func (e Employee) IsBirthday(date time.Time) bool {
    return e.Birthday.Year() != date.Year() &&
        e.Birthday.Month() == date.Month() &&
        e.Birthday.Day() == date.Day()
}

func (e Employee) Seniority(t time.Time) int {
    if t.Year() <= e.HireDate.Year() {
        return 0
    }

    seniority := t.Year() - e.HireDate.Year() - 1

    if t.Month() > e.HireDate.Month() || (t.Month() == e.HireDate.Month() && t.Day() >= e.HireDate.Day()) {
        seniority++
    }

    return seniority
}
