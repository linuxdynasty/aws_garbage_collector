package main

import (
	"log"
	"sync"

	"github.com/asdine/storm"
	"github.com/linuxdynasty/aws_garbage_collector/asgs"
	"github.com/linuxdynasty/aws_garbage_collector/datapipelines"
	"github.com/linuxdynasty/aws_garbage_collector/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/elasticache"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/rds"
	"github.com/linuxdynasty/aws_garbage_collector/redshift"
)

func ProcessEC2(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	e := ec2.DB(db)
	log.Printf("Begin Processing EC2 for region %s ", region)
	e.Process(region)
	log.Printf("Finished Processing EC2 for region %s ", region)
}

func ProcessElastiCache(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	r := elasticache.DB(db)
	log.Printf("Begin Processing ElastiCache for region %s ", region)
	r.Process(region)
	log.Printf("Finished Processing ElastiCache for region %s ", region)
}

func ProcessRedShift(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	r := redshift.DB(db)
	log.Printf("Begin Processing RedShift for region %s ", region)
	r.Process(region)
	log.Printf("Finished Processing Redshift for region %s ", region)
}

func ProcessRDS(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	r := rds.DB(db)
	log.Printf("Begin Processing RDS for region %s ", region)
	r.Process(region)
	log.Printf("Finished Processing RDS for region %s ", region)
}

func ProcessDataPipelines(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	dp := datapipelines.DB(db)
	log.Printf("Begin Processing Data Pipelines for region %s ", region)
	dp.Process(region)
	log.Printf("Finished Processing Data Pipelines for region %s ", region)
}

func ProcessIAM(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	awsIam := iam.DB(db)
	log.Printf("Begin Processing IAM")
	awsIam.Process(region)
	log.Printf("Finished Processing IAM")
}

func ProcessLaunchConfigurations(db *storm.DB, region string, wg *sync.WaitGroup) {
	defer wg.Done()
	lc := asgs.DB(db)
	log.Printf("Begin Processing Launch Configurations for region %s ", region)
	lc.Process(region)
	log.Printf("Finished Processing Launch Configurations for region %s ", region)
}
