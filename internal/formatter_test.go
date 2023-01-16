package internal

import "testing"

func TestSimpleFormatter_Format(t *testing.T) {
    type fields struct {
        Prefix string
        Suffix string
    }
    type args struct {
        names []Employee
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   string
    }{
        {fields: fields{Prefix: "Prefix", Suffix: "Suffix"}, args: args{[]Employee{{FirstName: "John", LastName: "Doe"}}}, want: "Prefix John Doe Suffix"},
        {fields: fields{Prefix: "Prefix", Suffix: "Suffix"}, args: args{[]Employee{{FirstName: "John", LastName: "Doe"}, {FirstName: "Jane", LastName: "Doe"}, {FirstName: "Will", LastName: "Smith"}}}, want: "Prefix Jane Doe, John Doe, Will Smith Suffix"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f := SimpleFormatter{
                Prefix: tt.fields.Prefix,
                Suffix: tt.fields.Suffix,
            }
            if got := f.Format(tt.args.names); got != tt.want {
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}
