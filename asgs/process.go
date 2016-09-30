package asgs

import (
	"log"
	"sync"
)

func (l *LC) Process(region string) {
	var wg sync.WaitGroup
	wg.Add(1)
	log.Printf("Getting Launch Configurations in region %s", region)
	go l.LaunchConfigurations(region, &wg)
	wg.Wait()
}
