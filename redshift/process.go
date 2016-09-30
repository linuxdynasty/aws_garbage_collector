package redshift

import (
	"log"
	"sync"
)

func (r *RedShift) Process(region string) error {
	var wg sync.WaitGroup
	wg.Add(1)

	log.Printf("Getting RedShift Security Groups in region %s", region)
	go r.RedShift(region, &wg)
	wg.Wait()
	return nil
}
