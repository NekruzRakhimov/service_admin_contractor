package service

import (
	"context"
	"service_admin_contractor/application/cerrors"
	"service_admin_contractor/domain/model"
	"service_admin_contractor/domain/repository"
	"time"
)

type ContractorService interface {
	FindContractors(ctx context.Context, params model.ContractorSearchParameters) ([]model.Contractor, int64, error)
	GetContractor(ctx context.Context, id int64) (model.Contractor, error)
	CreateContractor(ctx context.Context, contractor *model.Contractor) error
	UpdateContractor(ctx context.Context, id int64, contractor *model.Contractor) error
	DeleteContractor(id int64) error

	CreateContractorEmployee(ctx context.Context, contractorId int64, employee *model.Employee) error
	UpdateContractorEmployee(ctx context.Context, id int64, employee *model.Employee) error
	DeleteContractorEmployee(id int64) error
}

type contractorService struct {
	cr repository.ContractorRepository
}

func NewContractorService(cr repository.ContractorRepository) ContractorService {
	return &contractorService{cr}
}

func (cs *contractorService) FindContractors(ctx context.Context,
	params model.ContractorSearchParameters) ([]model.Contractor, int64, error) {
	result, total, err := cs.cr.FindContractors(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (cs *contractorService) GetContractor(ctx context.Context, id int64) (model.Contractor, error) {
	res, err := cs.cr.GetContractor(ctx, id)
	if err != nil {
		return model.Contractor{}, cerrors.ErrCouldNotGetContractorById(err, id)
	}

	return res, nil
}

func (cs *contractorService) CreateContractor(ctx context.Context, contractor *model.Contractor) error {
	tx, err := cs.cr.WithTransaction(ctx)
	if err != nil {
		return err
	}

	// Create Contractor
	err = cs.cr.CreateContractor(ctx, tx, contractor)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	return nil
}

func (cs *contractorService) UpdateContractor(ctx context.Context, id int64, contractor *model.Contractor) error {
	tx, err := cs.cr.WithTransaction(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	if contractor.Status == model.ContractorStatusBlock {
		blockDate := time.Now().UTC()
		contractor.BlockDate = &blockDate
	} else {
		contractor.BlockDate = nil
	}

	err = cs.cr.UpdateContractorData(ctx, tx, id, contractor)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	return nil
}

func (cs *contractorService) DeleteContractor(id int64) error {
	return cs.cr.DeleteContractor(id)
}

func (cs *contractorService) CreateContractorEmployee(ctx context.Context, contractorId int64,
	employee *model.Employee) error {
	tx, err := cs.cr.WithTransaction(ctx)
	if err != nil {
		return err
	}

	// Create Contractor employee
	err = cs.cr.CreateContractorEmployee(ctx, tx, contractorId, employee)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	return nil
}

func (cs *contractorService) UpdateContractorEmployee(ctx context.Context, id int64, employee *model.Employee) error {
	tx, err := cs.cr.WithTransaction(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	if employee.Status == model.EmployeeStatusBlock {
		blockDate := time.Now().UTC()
		employee.BlockDate = &blockDate
	}

	err = cs.cr.UpdateContractorEmployeeData(ctx, tx, id, employee)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return err
	}

	return nil
}

func (cs *contractorService) DeleteContractorEmployee(id int64) error {
	return cs.cr.DeleteContractorEmployee(id)
}
