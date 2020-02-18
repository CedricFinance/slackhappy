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
		Birthday string
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
			fields: fields{Birthday: "01-02"},
			args:   args{date: MustParseTime("Jan 2 2006", "Jan 2 2020")},
			want:   true,
		},
		{
			fields: fields{Birthday: "01-02"},
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
