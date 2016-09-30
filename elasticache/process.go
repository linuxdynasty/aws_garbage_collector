package elasticache

import (
	"log"
	"sync"
)

func (e *ElastiCache) Process(region string) error {
	var wg sync.WaitGroup
	wg.Add(1)

	log.Printf("Getting ElastiCache instances and clusters in region %s", region)
	go e.ElastiCache(region, &wg)
	wg.Wait()
	return nil
}
