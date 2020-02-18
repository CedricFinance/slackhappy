package internal

import (
	"context"
	"github.com/CedricFinance/slackhappy/bamboohr"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeClient bamboohr.CustomReportResponse

func (f FakeClient) CustomReport(ctx context.Context, title string, fields []string) (*bamboohr.CustomReportResponse, error) {
	r := bamboohr.CustomReportResponse(f)
	return &r, nil
}

func TestBambooRepository_List(t *testing.T) {
	type fields struct {
		Client *bamboohr.Client
	}

	b := &BambooRepository{
		Client: FakeClient(bamboohr.CustomReportResponse{
			Employees: []map[string]string{
				{
					"id":        "1",
					"status":    "Active",
					"hireDate":  "2010-10-14",
					"birthday":  "01-20",
					"firstName": "John",
					"lastName":  "Doe",
				},
				{
					"id":        "2",
					"status":    "Active",
					"hireDate":  "2010-03-17",
					"birthday":  "01-21",
					"firstName": "Alice",
					"lastName":  "Rabbit",
				},
				{
					"id":        "3",
					"status":    "Inactive",
					"hireDate":  "2016-09-02",
					"birthday":  "01-22",
					"firstName": "Bod",
					"lastName":  "Doe",
				},
				{
					"id":        "4",
					"status":    "Active",
					"hireDate":  "2016-09-03",
					"birthday":  "01-23",
					"firstName": "Marc",
					"lastName":  "Vador",
				},
				{
					"id":        "5",
					"status":    "Active",
					"hireDate":  "0000-00-00",
					"birthday":  "01-23",
					"firstName": "Marc",
					"lastName":  "Vador",
				},
			},
		}),
	}

	employees := b.List(context.Background())

	assert.Equal(t, 3, len(employees))
	assert.Equal(t, MustParseTime("2006-01-02", "2010-10-14"), employees[0].HireDate)
	assert.Equal(t, MustParseTime("2006-01-02", "2010-03-17"), employees[1].HireDate)
	assert.Equal(t, MustParseTime("2006-01-02", "2016-09-03"), employees[2].HireDate)
	assert.Equal(t, "01-20", employees[0].Birthday)
	assert.Equal(t, "01-21", employees[1].Birthday)
	assert.Equal(t, "01-23", employees[2].Birthday)
	assert.Equal(t, "John", employees[0].FirstName)
	assert.Equal(t, "Alice", employees[1].FirstName)
	assert.Equal(t, "Marc", employees[2].FirstName)
	assert.Equal(t, "Doe", employees[0].LastName)
	assert.Equal(t, "Rabbit", employees[1].LastName)
	assert.Equal(t, "Vador", employees[2].LastName)
}
