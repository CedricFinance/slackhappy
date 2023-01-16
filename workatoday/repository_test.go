package workatoday

import (
    "context"
    "github.com/CedricFinance/slackhappy/bamboohr"
    "github.com/stretchr/testify/assert"
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

type FakeClient WorkersResponse

func (f FakeClient) GetWorkersContext(ctx context.Context, requestType string) (*WorkersResponse, error) {
    r := WorkersResponse(f)
    return &r, nil
}

func TestBambooRepository_List(t *testing.T) {
    type fields struct {
        Client *bamboohr.Client
    }

    b := &EmployeeService{
        Client: FakeClient(
            WorkersResponse{
                Array: []Worker{
                    {
                        HireDate:  "2010-10-14",
                        Birthday:  "1980-01-20",
                        FirstName: "John",
                        LastName:  "Doe",
                    },
                    {
                        HireDate:  "2010-03-17",
                        Birthday:  "1981-01-21",
                        FirstName: "Alice",
                        LastName:  "Rabbit",
                    },
                    {
                        HireDate:  "2016-09-03",
                        Birthday:  "1982-01-23",
                        FirstName: "Marc",
                        LastName:  "Vador",
                    },
                },
            }),
    }

    employees, _ := b.ListContext(context.Background())

    assert.Equal(t, 3, len(employees))
    assert.Equal(t, MustParseTime("2006-01-02", "2010-10-14"), employees[0].HireDate)
    assert.Equal(t, MustParseTime("2006-01-02", "2010-03-17"), employees[1].HireDate)
    assert.Equal(t, MustParseTime("2006-01-02", "2016-09-03"), employees[2].HireDate)
    assert.Equal(t, MustParseTime("2006-01-02", "1980-01-20"), employees[0].Birthday)
    assert.Equal(t, MustParseTime("2006-01-02", "1981-01-21"), employees[1].Birthday)
    assert.Equal(t, MustParseTime("2006-01-02", "1982-01-23"), employees[2].Birthday)
    assert.Equal(t, "John", employees[0].FirstName)
    assert.Equal(t, "Alice", employees[1].FirstName)
    assert.Equal(t, "Marc", employees[2].FirstName)
    assert.Equal(t, "Doe", employees[0].LastName)
    assert.Equal(t, "Rabbit", employees[1].LastName)
    assert.Equal(t, "Vador", employees[2].LastName)
}
