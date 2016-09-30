package rds

import (
	"log"
	"sync"
)

func (r *RDS) Process(region string) error {
	var wg sync.WaitGroup
	wg.Add(1)

	log.Printf("Getting RDS Security Groups in region %s", region)
	go r.RDS(region, &wg)
	wg.Wait()
	return nil
}
