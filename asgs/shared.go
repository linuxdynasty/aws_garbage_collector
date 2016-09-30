package asgs

import "github.com/asdine/storm"

type LC struct {
	DB        *storm.DB
	ASGBucket storm.Node
}

func DB(db *storm.DB) LC {
	lc := LC{
		DB:        db,
		ASGBucket: db.From("LaunchConfiguration"),
	}
	return lc
}
