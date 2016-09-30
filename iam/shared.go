package iam

import "github.com/asdine/storm"

type IAM struct {
	DB *storm.DB
}

func DB(db *storm.DB) IAM {
	iam := IAM{
		DB: db,
	}
	return iam
}
