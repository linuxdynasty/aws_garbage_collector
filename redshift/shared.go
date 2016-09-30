package redshift

import "github.com/asdine/storm"

type RedShift struct {
	DB *storm.DB
}

func DB(db *storm.DB) RedShift {
	redshift := RedShift{
		DB: db,
	}
	return redshift
}
