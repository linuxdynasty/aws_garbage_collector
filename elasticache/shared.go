package elasticache

import "github.com/asdine/storm"

type ElastiCache struct {
	DB *storm.DB
}

func DB(db *storm.DB) ElastiCache {
	elasticache := ElastiCache{
		DB: db,
	}
	return elasticache
}
