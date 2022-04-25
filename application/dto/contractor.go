package dto

import (
	"github.com/go-playground/validator/v10"
	"net/url"
	"service_admin_contractor/domain/model"
	"time"
)

type ContractorDto struct {
	Id            int64         `json:"id"`
	Resident      bool          `json:"resident"`
	Bin           *string       `json:"bin"`
	Name          *string       `json:"name" validate:"required"`
	Email         string        `json:"email" validate:"required"`
	AgentName     *string        `json:"agentName" validate:"required"`
	AgentPassword string        `json:"agentPassword"`
	AgentPosition *string        `json:"agentPosition" validate:"required"`
	BlockDate     *time.Time    `json:"blockDate"`
	Status        string        `json:"status"`
	Employees     []EmployeeDto `json:"employees"`
}

func (dto ContractorDto) StructLevelValidation(sl validator.StructLevel) {
	if dto.Resident && (dto.Bin == nil || *dto.Bin == "") {
		sl.ReportError(dto.Bin, "bin", "Bin", "required", "")
	}
}

type EmployeeDto struct {
	Id        int64      `json:"id"`
	Email     string     `json:"email" validate:"required"`
	FullName  string     `json:"fullName"`
	Position  string     `json:"position" validate:"required"`
	BlockDate *time.Time `json:"blockDate"`
	Status    string     `json:"status"`
}

func ConvertContractors(list []model.Contractor) []interface{} {
	result := make([]interface{}, len(list))

	for i := range list {
		result[i] = ConvertContractor(list[i])
	}

	return result
}

func ConvertContractor(c model.Contractor) ContractorDto {
	employees := make([]EmployeeDto, 0)

	for _, e := range c.Employees {
		employees = append(employees, ConvertContractorEmployee(e))
	}

	return ContractorDto{
		Id:            c.Id,
		Resident:      c.Resident,
		Bin:           c.Bin,
		Name:          c.Name,
		Email:         c.Email,
		BlockDate:     c.BlockDate,
		Status:        string(c.Status),
		Employees:     employees,
		AgentName:     c.AgentName,
		AgentPassword: c.AgentPassword,
		AgentPosition: c.AgentPosition,
	}
}

func ConvertContractorEmployee(e model.Employee) EmployeeDto {
	return EmployeeDto{
		Id:        e.Id,
		Email:     e.Email,
		FullName:  e.FullName,
		Position:  e.Position,
		BlockDate: e.BlockDate,
		Status:    string(e.Status),
	}
}

func ParseContractorSearchParameters(values url.Values) (*model.ContractorSearchParameters, error) {
	pagination, err := ParsePagination(values)
	if err != nil {
		return nil, err
	}

	statusFilter := parseContractorStatusFilter(values)

	return &model.ContractorSearchParameters{
		Pagination: *pagination,
		Bin:        ParseStringFilter(values, "bin"),
		Name:       ParseStringFilter(values, "name"),
		Email:      ParseStringFilter(values, "email"),
		Status:     statusFilter,
	}, nil
}

func parseContractorStatusFilter(values url.Values) *model.ContractorStatus {
	filter := ParseStringFilter(values, "status")
	if filter == nil {
		return nil
	}

	result := model.ContractorStatus(*filter)

	switch result {
	case model.ContractorStatusActive,
		model.ContractorStatusBlock:
		return &result
	default:
		return nil
	}
}

func ConvertContractorDtoToEntity(dto *ContractorDto) *model.Contractor {
	return &model.Contractor{
		Resident:      dto.Resident,
		Bin:           dto.Bin,
		Name:          dto.Name,
		Email:         dto.Email,
		Status:        model.ContractorStatus(dto.Status),
		AgentName:     dto.AgentName,
		AgentPosition: dto.AgentPosition,
		AgentPassword: dto.AgentPassword,
	}
}

func ConvertEmployeeDtoToEntity(dto *EmployeeDto) *model.Employee {
	return &model.Employee{
		Email:    dto.Email,
		FullName: dto.FullName,
		Position: dto.Position,
		Status:   model.EmployeeStatus(dto.Status),
	}
}
