package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/sethvargo/go-password/password"
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

	GeneratePassword() (string, error)
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
	if len(contractor.AgentPassword) < 12 {
		return errors.New(fmt.Sprintf("не указан пароль или не соответсвует длина пароля"))
	}

	existingContractors, _, err := cs.FindContractors(ctx, model.ContractorSearchParameters{
		Pagination: *model.NewMaxPagination(),
		Email:      &contractor.Email})
	if err != nil {
		return cerrors.ErrCouldNotCreateContractor(err, " - поиск по базе возвратил ошибку")
	}

	if len(existingContractors) > 0 {
		errText := fmt.Sprintf("В базе уже есть email %s", contractor.Email)
		return cerrors.ErrCouldNotCreateContractor(errors.New(errText), errText)
	}

	tx, err := cs.cr.WithTransaction(ctx)
	if err != nil {
		return err
	}

	// Create Contractor
	if err = cs.cr.CreateContractor(ctx, tx, contractor); err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotCreateContractor(err, " - основные данные не записались в базу")
	}

	// Create Credentials
	if err = cs.createCredentials(ctx, tx, contractor); err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotCreateContractor(err, " - учетные данные не записались в базу")
	}

	err = tx.Commit(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotCreateContractor(err, " - не зафексирована в базе")
	}

	return nil
}

func (cs *contractorService) createCredentials(ctx context.Context, tx pgx.Tx, contractor *model.Contractor) error {
	var err error
	credentials := model.Credentials{
		ContractorId: &contractor.Id,
		Password:     contractor.AgentPassword,
	}

	if credentials.Password, err = credentials.GenerateHashPassword(); err != nil {
		return err
	}

	if err = cs.cr.CreateCredentials(ctx, tx, credentials); err != nil {
		return err
	}

	return err
}

func (cs *contractorService) UpdateContractor(ctx context.Context, id int64, contractor *model.Contractor) error {
	existingContractors, _, err := cs.FindContractors(ctx, model.ContractorSearchParameters{
		Pagination: *model.NewMaxPagination(),
		Email:      &contractor.Email})
	if err != nil {
		return cerrors.ErrCouldNotUpdateContractor(err, " - поиск по базе возвратил ошибку")
	}

	for _, c := range existingContractors {
		if c.Id != id {
			errText := fmt.Sprintf("В базе уже есть email %s", contractor.Email)
			return cerrors.ErrCouldNotUpdateContractor(errors.New(errText), errText)
		}
	}

	tx, err := cs.cr.WithTransaction(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotUpdateContractor(err, " - нет открылся транзакция")
	}

	if contractor.Status == model.ContractorStatusBlock {
		blockDate := time.Now().UTC()
		contractor.BlockDate = &blockDate
	} else {
		contractor.BlockDate = nil
		contractor.Status = model.ContractorStatusActive
	}

	if err = cs.cr.UpdateContractorData(ctx, tx, id, contractor); err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotUpdateContractor(err, " - основные данные не обновились")
	}

	if err = cs.updateContractorCredentials(ctx, tx, id, contractor); err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotUpdateContractor(err, " - данные по паролю не обновились")
	}

	err = tx.Commit(ctx)
	if err != nil {
		cs.cr.RollbackQuietly(tx, ctx)
		return cerrors.ErrCouldNotUpdateContractor(err, " - данные не обновились")
	}

	return nil
}

func (cs *contractorService) updateContractorCredentials(ctx context.Context, tx pgx.Tx, id int64,
	contractor *model.Contractor) error {
	var err error
	if len(contractor.AgentPassword) >= 12 {
		credentials := model.Credentials{
			ContractorId: &id,
			Password:     contractor.AgentPassword,
		}

		if credentials.Password, err = credentials.GenerateHashPassword(); err != nil {
			return err
		}

		if err = cs.cr.UpdateContractorCredentials(ctx, tx, credentials); err != nil {
			return err
		}
	}

	return err
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

func (cs *contractorService) GeneratePassword() (string, error) {
	generator, err := password.NewGenerator(&password.GeneratorInput{
		Symbols: model.Symbols,
	})
	if err != nil {
		return "", err
	}

	generatedPassword, err := generator.Generate(12, 4, 6, false, false)
	if err != nil {
		return "", err
	}

	return generatedPassword, nil
}
