package internal

import (
	"context"
	"fmt"
	"github.com/CedricFinance/slackhappy/bamboohr"
	"time"
)

type EmployeeRepository interface {
	List(context.Context) []Employee
}

type BambooRepository struct {
	Client bamboohr.Client
}

func (b *BambooRepository) List(ctx context.Context) []Employee {
	report, _ := b.Client.CustomReport(ctx, "Report for birthdays/anniversaries", []string{"firstName", "lastName", "birthday", "hireDate", "status"})

	employees := make([]Employee, 0, len(report.Employees))

	for _, employee := range report.Employees {

		if employee["status"] != "Active" {
			continue
		}

		hireDate, err := time.Parse("2006-01-02", employee["hireDate"])
		if err != nil {
			fmt.Printf("invalid hireDate %q\n", employee["hireDate"])
			continue
		}

		employees = append(employees, Employee{
			FirstName: employee["firstName"],
			LastName:  employee["lastName"],
			Birthday:  employee["birthday"],
			HireDate:  hireDate,
		})
	}

	return employees
}
