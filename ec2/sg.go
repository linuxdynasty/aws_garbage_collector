package ec2

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func (e *EC2) sg(securityGroups []*ec2.SecurityGroup, region string) error {
	for _, val := range securityGroups {
		group := models.SecurityGroup{
			ID:          *val.GroupId,
			Name:        *val.GroupName,
			Description: *val.Description,
			Region:      region,
		}
		if err := e.DB.Save(&group); err != nil {
			log.Fatal(err)
		}
		for _, awstag := range val.Tags {
			if *awstag.Key != "Name" {
				tag := models.Tag{
					ResourceId: *val.GroupId,
					Key:        *awstag.Key,
					Value:      *awstag.Value,
					Region:     region,
				}
				if err := e.DB.Save(&tag); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return nil
}

func (e *EC2) sgSources(securityGroups []*ec2.SecurityGroup) error {
	sgMatch := regexp.MustCompile("^sg-")
	sDb := e.DB.From("SecurityGroup")
	for _, sg := range securityGroups {
		var sgs map[string]bool
		sgs = make(map[string]bool)
		ipPermissions := append(sg.IpPermissions, sg.IpPermissionsEgress...)
		for _, cidr := range ipPermissions {
			for _, groupPair := range cidr.UserIdGroupPairs {
				if ok := sgMatch.MatchString(*groupPair.GroupId); ok {
					sgs[*groupPair.GroupId] = true
				}
			}
		}
		for sgid := range sgs {
			group := &models.SecurityGroup{}
			if err := e.DB.One("ID", sgid, group); err != nil {
				log.Print(sgid, " - not found")
			} else {
				source := models.SourceSecurityGroup{
					GroupId:       *sg.GroupId,
					SourceGroupId: group.ID,
					Name:          group.Name,
				}
				if err := sDb.Save(&source); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return nil
}

func (e *EC2) securityGroups(client *ec2.EC2) error {
	params := &ec2.DescribeSecurityGroupsInput{}
	resp, err := client.DescribeSecurityGroups(params)

	if err != nil {
		return err
	}
	err = e.sg(resp.SecurityGroups, *client.Config.Region)
	if err != nil {
		return err
	}
	err = e.sgSources(resp.SecurityGroups)
	if err != nil {
		return err
	}
	return nil
}

func (e *EC2) SGDeleteByIds(groupIds []string, region string) []models.DeleteStatus {
	session := session.New(&aws.Config{Region: &region})
	ec2_svc := ec2.New(session)
	statuses := []models.DeleteStatus{}

	for _, groupId := range groupIds {
		var sg models.SecurityGroup
		params := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(groupId),
		}
		_, err := ec2_svc.DeleteSecurityGroup(params)
		status := models.DeleteStatus{
			ID: groupId,
		}
		if queryErr := e.DB.One("ID", groupId, &sg); queryErr == nil {
			status.Name = sg.Name
		}
		if err != nil {
			status.Deleted = false
			status.Message = err.Error()
		} else {
			status.Deleted = true
			status.Message = fmt.Sprintf("%s deleted successfully", groupId)
		}
		statuses = append(statuses, status)
	}
	return statuses
}

func (e *EC2) StoreSecurityGroups(region string) error {
	var err error
	session := session.New(&aws.Config{Region: &region})
	ec2_svc := ec2.New(session)

	if sgErr := e.securityGroups(ec2_svc); err != nil {
		err = fmt.Errorf("failed to get security groups: %v", sgErr)
		return err
	}
	return err
}
