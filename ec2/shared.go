package ec2

import "github.com/asdine/storm"

type EC2 struct {
	DB       *storm.DB
	LcBucket storm.Node
}

func DB(db *storm.DB) EC2 {
	ec2 := EC2{
		DB:       db,
		LcBucket: db.From("LaunchConfiguration"),
	}
	return ec2
}
