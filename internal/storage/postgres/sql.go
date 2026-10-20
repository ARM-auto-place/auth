package postgres

import (
	_ "embed"
)

//go:embed sql/insert_user.sql
var sqlInsertUser string

//go:embed sql/select_user_by_email.sql
var sqlSelectUserByEmail string
