package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"service_admin_contractor/domain/model"
	"service_admin_contractor/infrastructure/persistence/postgres"
)

type ContractorRepository interface {
	postgres.Transactional
	FindContractors(ctx context.Context, params model.ContractorSearchParameters) ([]model.Contractor, int64, error)
	GetContractor(ctx context.Context, id int64) (model.Contractor, error)
	CreateContractor(ctx context.Context, tx pgx.Tx, contractor *model.Contractor) error
	UpdateContractorData(ctx context.Context, tx pgx.Tx, contractorId int64, contractor *model.Contractor) error
	DeleteContractor(id int64) error

	CreateContractorEmployee(ctx context.Context, tx pgx.Tx, contractorId int64, employee *model.Employee) error
	UpdateContractorEmployeeData(ctx context.Context, tx pgx.Tx, employeeId int64, employee *model.Employee) error
	DeleteContractorEmployee(id int64) error
}
