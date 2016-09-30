package api

import (
	"github.com/asdine/storm/q"
	"github.com/linuxdynasty/aws_garbage_collector/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func getAllResources(e *ec2.EC2, groups []models.SecurityGroup) []models.SecurityGroupDetails {
	var groupDetails []models.SecurityGroupDetails

	for _, sg := range groups {
		groupDetail := models.SecurityGroupDetails{
			ID:            sg.ID,
			Name:          sg.Name,
			Description:   sg.Description,
			InUse:         sg.InUse,
			InUseBySGOnly: sg.InUseBySGOnly,
			Region:        sg.Region,
		}
		var ec2 []models.SecurityGroupResource
		var elasticache []models.SecurityGroupResource
		var elb []models.SecurityGroupResource
		var elbv2 []models.SecurityGroupResource
		var launchconfiguration []models.SecurityGroupResource
		var rds []models.SecurityGroupResource
		var redshift []models.SecurityGroupResource
		var sgs []models.SourceSecurityGroup
		var tags []models.Tag

		gDb := e.DB.From("SecurityGroup")
		sDb := e.DB.From("SecurityGroup")

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "EC2"),
		)).Find(&ec2)
		groupDetail.EC2 = ec2

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "ElastiCache"),
		)).Find(&elasticache)
		groupDetail.ElastiCache = elasticache

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "ELB"),
		)).Find(&elb)
		groupDetail.ELB = elb

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "ELBV2"),
		)).Find(&elbv2)
		groupDetail.ELBV2 = elbv2

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "Launch Configuration"),
		)).Find(&launchconfiguration)
		groupDetail.LaunchConfiguration = launchconfiguration

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "RDS"),
		)).Find(&rds)
		groupDetail.RDS = rds

		gDb.Select(q.And(
			q.Eq("GroupId", sg.ID),
			q.Eq("AWSType", "Redshift"),
		)).Find(&redshift)
		groupDetail.RedShift = redshift

		sDb.Find("GroupId", sg.ID, &sgs)
		groupDetail.SourceSecurityGroups = sgs

		e.DB.Find("ResourceId", sg.ID, &tags)
		groupDetail.Tags = tags

		groupDetails = append(groupDetails, groupDetail)
	}
	return groupDetails
}

func GetAllSecurityGroups(e *ec2.EC2, region, inUse, inUseBySgOnly string, all bool) *models.GetResponseApi {
	var groups []models.SecurityGroup
	var queries []q.Matcher
	if inUse != "" {
		queries = append(queries, q.Eq("InUse", inUse))
	}
	if inUseBySgOnly != "" {
		queries = append(queries, q.Eq("InUseBySGOnly", inUseBySgOnly))
	}
	if region != "" {
		queries = append(queries, q.Eq("Region", region))
		if all {
			e.DB.Find("Region", region, &groups)
		} else {
			e.DB.Select(q.And(
				queries...,
			)).Find(&groups)
		}
	} else {
		if all {
			e.DB.All(&groups)
		} else {
			e.DB.Select(q.And(
				queries...,
			)).Find(&groups)
		}
	}
	groupDetails := getAllResources(e, groups)

	apiResponse := models.GetResponseApi{
		Data:  groupDetails,
		Count: len(groupDetails),
	}
	return &apiResponse
}
