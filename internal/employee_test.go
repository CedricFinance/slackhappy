package internal

import (
    "log"
    "testing"
    "time"
)

func MustParseTime(layout, value string) time.Time {
    res, err := time.Parse(layout, value)
    if err != nil {
        log.Panic(err)
    }
    return res
}

func TestEmployee_IsAnniversary(t *testing.T) {
    type fields struct {
        HireDate time.Time
    }
    type args struct {
        date time.Time
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   bool
    }{
        {
            fields: fields{HireDate: MustParseTime("2006-01-02", "2013-01-02")},
            args:   args{date: MustParseTime("Jan 2 2006", "Jan 2 2020")},
            want:   true,
        },
        {
            fields: fields{HireDate: MustParseTime("2006-01-02", "2013-01-02")},
            args:   args{date: MustParseTime("Jan 2 2006", "Jan 3 2020")},
            want:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := Employee{
                HireDate: tt.fields.HireDate,
            }
            if got := e.IsAnniversary(tt.args.date); got != tt.want {
                t.Errorf("IsAnniversary() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestEmployee_IsBirthday(t *testing.T) {
    type fields struct {
        Birthday time.Time
    }
    type args struct {
        date time.Time
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   bool
    }{
        {
            fields: fields{Birthday: MustParseTime("Jan 2 2006", "Jan 2 1980")},
            args:   args{date: MustParseTime("Jan 2 2006", "Jan 2 2020")},
            want:   true,
        },
        {
            fields: fields{Birthday: MustParseTime("Jan 2 2006", "Jan 2 1980")},
            args:   args{date: MustParseTime("Jan 2 2006", "Jan 3 2020")},
            want:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := Employee{
                Birthday: tt.fields.Birthday,
            }
            if got := e.IsBirthday(tt.args.date); got != tt.want {
                t.Errorf("IsBirthday() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestEmployee_Seniority(t *testing.T) {
    type args struct {
        t time.Time
    }
    tests := []struct {
        name        string
        hireDate    date
        currentDate date
        want        int
    }{
        {"before hire date", date{year: 2020, month: time.July, day: 1}, date{year: 2010, month: time.January, day: 1}, 0},
        {"before 1 year", date{year: 2020, month: time.July, day: 1}, date{year: 2020, month: time.August, day: 1}, 0},
        {"before 1 year", date{year: 2020, month: time.July, day: 1}, date{year: 2021, month: time.June, day: 30}, 0},
        {"exactly 1 year", date{year: 2020, month: time.July, day: 1}, date{year: 2025, month: time.October, day: 10}, 5},
        {"several years", date{year: 2020, month: time.July, day: 1}, date{year: 2025, month: time.October, day: 10}, 5},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := Employee{
                HireDate: newDate(tt.hireDate),
            }
            if got := e.Seniority(newDate(tt.currentDate)); got != tt.want {
                t.Errorf("Seniority() = %v, want %v", got, tt.want)
            }
        })
    }
}

type date struct {
    year  int
    month time.Month
    day   int
}

func newDate(d date) time.Time {
    return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.UTC)
}
