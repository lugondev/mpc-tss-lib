package db

import "database/sql"

func NullString() sql.NullString {
	return sql.NullString{
		Valid: true,
	}
}
