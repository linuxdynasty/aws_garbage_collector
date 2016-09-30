package datapipelines

import "github.com/asdine/storm"

type DataPipeline struct {
	DB *storm.DB
}

func DB(db *storm.DB) DataPipeline {
	dp := DataPipeline{
		DB: db,
	}
	return dp
}
