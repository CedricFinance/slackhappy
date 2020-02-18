package internal

import (
	"fmt"
	"strings"
)

type Formatter interface {
	Format(names []string) string
}

type SimpleFormatter struct {
	Prefix, Suffix string
}

func (f SimpleFormatter) Format(names []string) string {
	joinedNames := strings.Join(names, ", ")
	return fmt.Sprintf("%s %s %s", f.Prefix, joinedNames, f.Suffix)
}
