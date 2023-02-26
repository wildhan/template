package repository

import "template/config/database"

type authRepo struct {
	dbConn *database.DbConnection
}

func NewAuthRepo(dbConn *database.DbConnection) AuthRepo {
	return &authRepo{dbConn}
}

type AuthRepo interface {
}
