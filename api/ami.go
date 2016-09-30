package api

import (
	"github.com/asdine/storm/q"
	"github.com/linuxdynasty/aws_garbage_collector/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func GetAllUnusedAMIIds(e *ec2.EC2, region string) []string {
	var amis []models.EC2Ami
	var amiIds []string
	e.DB.Select(q.And(
		q.Eq("InUseByLC", "false"),
		q.Eq("InUseByInstance", "false"),
		q.Eq("InUseByDataPipeline", "false"),
		q.Eq("Region", region),
	)).Find(&amis)

	for _, ami := range amis {
		amiIds = append(amiIds, ami.ID)
	}

	return amiIds
}

func getAllAmiAttachments(e *ec2.EC2, amis []models.EC2Ami) []models.EC2AmiDetails {
	var amiDetails []models.EC2AmiDetails
	for _, ami := range amis {
		amiDetail := models.EC2AmiDetails{
			ID:                  ami.ID,
			Name:                ami.Name,
			Description:         ami.Description,
			ImageLocation:       ami.ImageLocation,
			Region:              ami.Region,
			VirtualizationType:  ami.VirtualizationType,
			SnapshotIds:         ami.SnapshotIds,
			State:               ami.State,
			InUseByLC:           ami.InUseByLC,
			InUseByInstance:     ami.InUseByInstance,
			InUseByDataPipeline: ami.InUseByDataPipeline,
		}
		var tags []models.Tag
		var ec2Data []models.EC2Instance
		var lcData []models.LaunchConfiguration
		var dpData []models.PipeLine
		e.DB.Find("ImageId", ami.ID, &ec2Data)
		e.DB.Find("ImageId", ami.ID, &lcData)
		e.DB.Find("ImageId", ami.ID, &dpData)
		e.DB.Find("ResourceId", ami.ID, &tags)
		amiDetail.EC2 = ec2Data
		amiDetail.LC = lcData
		amiDetail.Tags = tags
		amiDetail.DataPipelines = dpData
		amiDetails = append(amiDetails, amiDetail)
	}
	return amiDetails
}

func GetAllAMIs(e *ec2.EC2, region, inUse, inUseByLC, inUseByInstance, inUseByDataPipeline string, all bool) *models.GetResponseApi {
	var amis []models.EC2Ami
	var queries []q.Matcher

	if inUseByLC != "" {
		queries = append(queries, q.Eq("InUseByLC", inUseByLC))
	}
	if inUseByInstance != "" {
		queries = append(queries, q.Eq("InUseByInstance", inUseByInstance))
	}
	if inUseByDataPipeline != "" {
		queries = append(queries, q.Eq("InUseByInstance", inUseByDataPipeline))
	}
	if region != "" {
		queries = append(queries, q.Eq("Region", region))
		if all {
			e.DB.Find("Region", region, &amis)
		} else if inUse != "" {
			e.DB.Select(q.And(
				q.Eq("InUse", inUse),
				q.Eq("Region", region),
			)).Find(&amis)
		} else {
			e.DB.Select(q.And(
				queries...,
			)).Find(&amis)
		}
	} else {
		if all {
			e.DB.All(&amis)
		} else if inUse != "" {
			e.DB.Find("InUse", inUse, &amis)
		} else {
			e.DB.Select(q.And(
				queries...,
			)).Find(&amis)
		}
	}

	amiDetails := getAllAmiAttachments(e, amis)

	apiResponse := models.GetResponseApi{
		Data:  amiDetails,
		Count: len(amiDetails),
	}
	return &apiResponse
}
