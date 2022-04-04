package repository

type BpmsUserRepository interface {
	FindUserRoles (login string) ([]string, error)
}
