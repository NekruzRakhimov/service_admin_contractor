package postgres

import "github.com/jackc/pgx/v4/pgxpool"

type BpmsUserRepository struct {
	db *pgxpool.Pool
}

func NewBpmsUserRepository(db *pgxpool.Pool) *BpmsUserRepository {
	return &BpmsUserRepository{db}
}

func (b *BpmsUserRepository) FindUserRoles(login string) ([]string, error) {
	return nil, nil
}
