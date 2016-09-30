package rds

import "github.com/asdine/storm"

type RDS struct {
	DB *storm.DB
}

func DB(db *storm.DB) RDS {
	rds := RDS{
		DB: db,
	}
	return rds
}
