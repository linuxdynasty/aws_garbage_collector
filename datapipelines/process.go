package datapipelines

import (
	"log"
	"sync"
)

func (p *DataPipeline) Process(region string) {
	var wg sync.WaitGroup
	wg.Add(1)
	log.Printf("Getting Pipelines in region %s", region)
	go p.FetchAndStorePipelines(region, &wg)
	wg.Wait()
}
