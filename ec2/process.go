package ec2

import (
	"log"
	"sync"

	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func (e *EC2) SGProcessInUse() {
	var groups []models.SecurityGroup
	e.DB.All(&groups)
	rDb := e.DB.From("SecurityGroup")
	sDb := e.DB.From("SecurityGroup")
	for _, sg := range groups {
		var resources []models.SecurityGroupResource
		var sgs []models.SourceSecurityGroup
		rDb.Find("GroupId", sg.ID, &resources)
		sDb.Find("GroupId", sg.ID, &sgs)
		if len(sgs) > 0 && len(resources) == 0 {
			sg.InUse = "false"
			sg.InUseBySGOnly = "true"
		} else if len(resources) > 0 {
			sg.InUseBySGOnly = "false"
			sg.InUse = "true"
		} else if len(sgs) == 0 && len(resources) == 0 {
			sg.InUse = "false"
			sg.InUseBySGOnly = "false"
		}
		e.DB.Update(&sg)
	}
}

func (e *EC2) EC2ProcessInUse() {
	var amis []models.EC2Ami
	e.DB.All(&amis)
	for _, ami := range amis {
		ami.InUseByLC = "false"
		ami.InUseByInstance = "false"
		ami.InUseByDataPipeline = "false"
		var ec2s []models.EC2Instance
		var lcs []models.LaunchConfiguration
		var dps []models.PipeLine
		e.DB.Find("ImageId", ami.ID, &ec2s)
		e.DB.Find("ImageId", ami.ID, &lcs)
		e.DB.Find("ImageId", ami.ID, &dps)
		if len(ec2s) > 0 {
			ami.InUseByInstance = "true"
		}
		if len(lcs) > 0 {
			ami.InUseByLC = "true"
		}
		if len(dps) > 0 {
			ami.InUseByDataPipeline = "true"
		}
		if ami.InUseByInstance == "false" && ami.InUseByLC == "false" && ami.InUseByDataPipeline == "false" {
			ami.InUse = "false"
		}
		e.DB.Update(&ami)
	}
}

func (e *EC2) Process(region string) {
	var wg sync.WaitGroup
	wg.Add(3)
	log.Printf("Getting EC2 in region %s", region)
	go e.FetchAndStoreEC2(region, &wg)
	go e.StoreAMIs(region, &wg)
	log.Printf("Getting EC2 SecurityGroups in region %s", region)
	if err := e.StoreSecurityGroups(region); err != nil {
		log.Fatal(err)
	}
	log.Printf("Getting ELBv1 and ELBv2 SecurityGroups in region %s", region)
	go e.StoreELBs(region, &wg)
	wg.Wait()
}
