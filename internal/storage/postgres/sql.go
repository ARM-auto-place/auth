package postgres

import (
	_ "embed"
)

//go:embed sql/insert_user.sql
var sqlInsertUser string
