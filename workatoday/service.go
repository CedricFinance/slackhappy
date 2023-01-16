package workatoday

import (
    "context"
    "fmt"
    "github.com/CedricFinance/slackhappy/internal"
    "net/http"
    "time"
)

type EmployeeService struct {
    Client WorkatodayClient
}

type Config struct {
    BaseURL string
    ApiKey  string
}

func NewEmployeeService(c *Config, client *http.Client) *EmployeeService {
    employeesService := EmployeeService{
        Client: &Client{
            BaseURL: c.BaseURL,
            APIKey:  c.ApiKey,
            Client:  client,
        },
    }
    return &employeesService
}

func (s *EmployeeService) ListContext(ctx context.Context) ([]internal.Employee, error) {
    report, err := s.Client.GetWorkersContext(ctx, "blabla_happy")
    if err != nil {
        return nil, fmt.Errorf("failed to get the custom report: %v", err)
    }

    var result []internal.Employee

    for _, employee := range report.Records {
        if employee.Active != "1" {
            continue
        }

        hireDate, _ := time.Parse("2006-01-02", employee.HireDate)
        birthday, _ := time.Parse("2006-01-02", employee.Birthday)

        result = append(result, internal.Employee{
            FirstName: employee.FirstName,
            LastName:  employee.LastName,
            Birthday:  birthday,
            HireDate:  hireDate,
        })
    }

    return result, nil
}
