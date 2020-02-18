package internal

import "testing"

func TestSimpleFormatter_Format(t *testing.T) {
	type fields struct {
		Prefix string
		Suffix string
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{fields: fields{Prefix: "Prefix", Suffix: "Suffix"}, args: args{[]string{"test"}}, want: "Prefix test Suffix"},
		{fields: fields{Prefix: "Prefix", Suffix: "Suffix"}, args: args{[]string{"test", "test2", "test3"}}, want: "Prefix test, test2, test3 Suffix"},
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
