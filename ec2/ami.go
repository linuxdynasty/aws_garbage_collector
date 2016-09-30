package ec2

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func (e *EC2) storeAMIs(client *ec2.EC2) {
	owner := "self"
	params := &ec2.DescribeImagesInput{
		Owners: []*string{&owner},
	}
	resp, err := client.DescribeImages(params)
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range resp.Images {
		for _, tag := range image.Tags {
			awstag := models.Tag{
				ResourceId: *image.ImageId,
				Key:        *tag.Key,
				Value:      *tag.Value,
				Region:     *client.Config.Region,
			}
			if err := e.DB.Save(&awstag); err != nil {
				log.Fatal(err)
			}
		}
		ami := models.EC2Ami{
			ID:                 *image.ImageId,
			Name:               *image.Name,
			Region:             *client.Config.Region,
			VirtualizationType: *image.VirtualizationType,
			State:              *image.State,
		}
		if image.Description != nil {
			ami.Description = *image.Description
		}
		if image.ImageLocation != nil {
			ami.ImageLocation = *image.ImageLocation
		}
		var snapIds []string
		for _, blockDevice := range image.BlockDeviceMappings {
			if blockDevice.Ebs != nil {
				snapIds = append(snapIds, *blockDevice.Ebs.SnapshotId)
			}
		}
		ami.SnapshotIds = snapIds
		if err := e.DB.Save(&ami); err != nil {
			log.Fatal(err)
		}
	}
}

func (e *EC2) DeleteSnapshotByIds(snapIds []string, region string) ([]string, []string) {
	var deletedIds []string
	var failedIds []string

	session := session.New(&aws.Config{Region: &region})
	ec2_svc := ec2.New(session)

	for _, snapId := range snapIds {
		params := &ec2.DeleteSnapshotInput{
			SnapshotId: aws.String(snapId),
		}
		_, err := ec2_svc.DeleteSnapshot(params)
		if err != nil {
			failedIds = append(failedIds, snapId)
		} else {
			deletedIds = append(deletedIds, snapId)
		}
	}
	return deletedIds, failedIds
}

func (e *EC2) DeleteByIds(amiIds []string, region string) []models.DeleteStatus {
	session := session.New(&aws.Config{Region: &region})
	ec2_svc := ec2.New(session)
	statuses := []models.DeleteStatus{}

	for _, amiId := range amiIds {
		var ami models.EC2Ami
		params := &ec2.DeregisterImageInput{
			ImageId: aws.String(amiId),
		}
		_, err := ec2_svc.DeregisterImage(params)
		status := models.DeleteStatus{
			ID: amiId,
		}
		if queryErr := e.DB.One("ID", amiId, &ami); queryErr == nil {
			status.Name = ami.Name
		}
		if err != nil {
			status.Deleted = false
			status.Message = err.Error()
		} else {
			status.Deleted = true
			snapsDeleted, snapsFailed := e.DeleteSnapshotByIds(ami.SnapshotIds, region)
			if len(snapsFailed) > 0 {
				status.Message = fmt.Sprintf("Failed to delete the following snapshot ids: %v for AMI: %s", snapsFailed, amiId)
			} else {
				status.Message = fmt.Sprintf("AMI: %s deleted successfully with it's snapshots: %v", amiId, snapsDeleted)
			}
			e.DB.DeleteStruct(&ami)
		}
		statuses = append(statuses, status)
	}
	return statuses
}

func (e *EC2) StoreAMIs(region string, wg *sync.WaitGroup) {
	defer wg.Done()
	session := session.New(&aws.Config{Region: &region})
	svc := ec2.New(session)
	e.storeAMIs(svc)
}
