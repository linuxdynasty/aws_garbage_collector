package iam

import (
	"log"
	"sync"

	"github.com/asdine/storm/q"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (i *IAM) InstanceProfileProcessInUse() {
	var instanceProfiles []models.IAMInstanceProfile
	i.DB.All(&instanceProfiles)
	rDb := i.DB.From("IAMInstanceProfile")
	for _, instanceProfile := range instanceProfiles {
		var instances []models.IAMResource
		var lcs []models.IAMResource
		rDb.Select(q.And(
			q.Eq("Type", shared.Instance),
			q.Eq("ARN", instanceProfile.ARN),
		)).Find(&instances)
		rDb.Select(q.And(
			q.Eq("Type", shared.LaunchConfiguration),
			q.Eq("ARN", instanceProfile.ARN),
		)).Find(&lcs)

		if len(instances) > 0 {
			instanceProfile.InUseByInstances = "true"
			instanceProfile.InUseByLCs = "true"
		}
		if len(lcs) > 0 {
			instanceProfile.InUseByLCs = "true"
		}
		i.DB.Update(&instanceProfile)
	}
}

func (i *IAM) Process(region string) {
	var wg sync.WaitGroup
	wg.Add(1)
	log.Printf("Getting IAM Policies")
	go i.FetchAndStoreIAM(region, &wg)
	wg.Wait()
}
