package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"service_admin_contractor/domain/model"
)

type ContractorRepository struct {
	db *pgxpool.Pool
}

func NewContractorRepository(db *pgxpool.Pool) *ContractorRepository {
	return &ContractorRepository{db}
}

func (c *ContractorRepository) RollbackQuietly(tx pgx.Tx, ctx context.Context) {
	err := tx.Rollback(ctx)
	if err != nil {
		log.Warn(err)
	}
}

func (c *ContractorRepository) WithTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *ContractorRepository) FindContractors(ctx context.Context,
	params model.ContractorSearchParameters) ([]model.Contractor, int64, error) {
	args := model.NamedArguments{}
	queryTotal := `select count(*)`
	querySelect := `select c.id, c.resident, c.bin, c.name, c.email, c.block_date, c.status,
							(
								SELECT
									JSON_AGG(e.*)
								FROM contractors_contractor_employee e WHERE e.contractor_id = c.id and e.is_delete = false
							) as employees`
	queryFrom := ` from contractors_contractor c`
	filters := ` where 1=1 and c.is_delete = false`

	AppendEqualsFilter(&filters, args, "c.bin", params.Bin)
	AppendStringLikeFilter(&filters, args, "c.name", params.Name, "%s%%")
	AppendStringLikeFilter(&filters, args, "c.email", params.Email, "%s%%")
	AppendEqualsFilter(&filters, args, "c.status", params.Status)

	var total int64
	_, err := QueryWithMap(c.db, ctx, queryTotal+queryFrom+filters, args).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []model.Contractor{}, 0, nil
	}

	paginatedFilters := filters + ` order by c.id desc`
	AppendPagination(&paginatedFilters, args, params.Pagination)

	result, err := QueryWithMap(c.db, ctx, querySelect+queryFrom+paginatedFilters, args).ReadAll(model.Contractor{})
	if err != nil {
		return nil, 0, err
	}

	return result.([]model.Contractor), total, nil
}

func (c *ContractorRepository) GetContractor(ctx context.Context, id int64) (model.Contractor, error) {
	args := make(model.NamedArguments)
	args["id"] = id
	query := `SELECT c.id, c.resident, c.bin, c.name, c.email, c.block_date, c.status,
					(
						SELECT
							JSON_AGG(e.*)
						FROM contractors_contractor_employee e WHERE e.contractor_id = c.id and e.is_delete = false
					) as employees
				FROM contractors_contractor c
						 where c.id = :id and c.is_delete = false`
	res, err := QueryWithMap(c.db, ctx, query, args).Read(model.Contractor{})
	if err != nil {
		return model.Contractor{}, err
	}

	return c.unwrapContractorSlice(res), nil
}

func (c *ContractorRepository) unwrapContractorSlice(res interface{}) model.Contractor {
	if res == nil {
		return model.Contractor{}
	} else {
		contractor := res.(*model.Contractor)
		return *contractor
	}
}

func (c *ContractorRepository) CreateContractor(ctx context.Context, tx pgx.Tx, contractor *model.Contractor) error {
	query := `INSERT INTO contractors_contractor (
					 resident, bin, name, email, status,agent_name,agent_position
				) VALUES (
					:resident, :bin, :name, :email, :status,:agent_name,:agent_position
				) RETURNING id`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"resident":       contractor.Resident,
		"bin":            contractor.Bin,
		"name":           contractor.Name,
		"email":          contractor.Email,
		"status":         model.ContractorStatusActive,
		"agent_name":     contractor.AgentName,
		"agent_position": contractor.AgentPosition,
	})
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, finalQuery, queryArgs...).Scan(&contractor.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractorRepository) UpdateContractorData(ctx context.Context, tx pgx.Tx, contractorId int64,
	contractor *model.Contractor) error {
	query := `UPDATE contractors_contractor 
				SET
					resident = 		:resident, 
					bin = 			:bin, 
					name = 			:name, 
					email = 		:email, 
					block_date =	:block_date,
					status = 		:status
				WHERE ID = :id_value`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"resident":   contractor.Resident,
		"bin":        contractor.Bin,
		"name":       contractor.Name,
		"email":      contractor.Email,
		"block_date": contractor.BlockDate,
		"status":     contractor.Status,
		"id_value":   contractorId,
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, finalQuery, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractorRepository) DeleteContractor(id int64) error {
	query := `update contractors_contractor 
				set is_delete = true where id = :id`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	_, err = c.db.Exec(context.Background(), finalQuery, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractorRepository) CreateContractorEmployee(ctx context.Context, tx pgx.Tx, contractorId int64,
	employee *model.Employee) error {
	query := `INSERT INTO contractors_contractor_employee (
					 contractor_id, email, full_name,  position
				) VALUES (
					:contractor_id, :email, :full_name, :position
				) RETURNING id`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"contractor_id": contractorId,
		"email":         employee.Email,
		"full_name":     employee.FullName,
		"position":      employee.Position,
	})
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, finalQuery, queryArgs...).Scan(&employee.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContractorRepository) UpdateContractorEmployeeData(ctx context.Context, tx pgx.Tx, employeeId int64,
	employee *model.Employee) error {
	query := `UPDATE contractors_contractor_employee 
				SET
					email = 		:email,
					full_name = 	:full_name,
					position = 		:position, 
					block_date =	:block_date,
					status = 		:status
				WHERE ID = :id_value`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"email":      employee.Email,
		"full_name":  employee.FullName,
		"position":   employee.Position,
		"block_date": employee.BlockDate,
		"status":     employee.Status,
		"id_value":   employeeId,
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, finalQuery, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractorRepository) DeleteContractorEmployee(id int64) error {
	query := `update contractors_contractor_employee 
				set is_delete = true where id = :id`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	_, err = c.db.Exec(context.Background(), finalQuery, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractorRepository) CreateCredentials(ctx context.Context, tx pgx.Tx, credentials model.Credentials) error {
	query := `INSERT INTO contractors_credentials (
					 contractor_id, employee_id, password
				) VALUES (
					:contractor_id, :employee_id, :password
				) RETURNING id`

	finalQuery, queryArgs, err := InlineNamedPlaceholders(query, map[string]interface{}{
		"contractor_id": credentials.ContractorId,
		"employee_id":   credentials.EmployeeId,
		"password":      credentials.Password,
	})
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, finalQuery, queryArgs...).Scan(&credentials.Id)
	if err != nil {
		return err
	}
	return nil
}
