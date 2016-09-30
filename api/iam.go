package api

import (
	"github.com/asdine/storm/q"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func GetAllUnusedPolicyIds(i *iam.IAM, region string) []string {
	var policies []models.IAMManagedPolicy
	var policyIds []string
	i.DB.Select(q.And(
		q.Eq("AttachmentCount", 0),
		q.Eq("Region", region),
	)).Find(&policies)

	for _, policy := range policies {
		policyIds = append(policyIds, policy.ID)
	}

	return policyIds
}

func getAllPolicyAttachments(i *iam.IAM, policies []models.IAMManagedPolicy) []models.IAMManagedPolicyDetails {
	var policyDetails []models.IAMManagedPolicyDetails
	for _, policy := range policies {
		policyDetail := models.IAMManagedPolicyDetails{
			ARN:              policy.ARN,
			AttachementCount: policy.AttachementCount,
			ID:               policy.ID,
			InUse:            policy.InUse,
			InUseByUsers:     policy.InUseByUsers,
			InUseByRoles:     policy.InUseByRoles,
			InUseByGroups:    policy.InUseByGroups,
			Name:             policy.Name,
			VersionId:        policy.VersionId,
			Policy:           policy.Policy,
			Region:           policy.Region,
		}
		rDB := i.DB.From("IAMManagedPolicy")
		rDB.Select(q.And(
			q.Eq("Type", "Groups"),
			q.Eq("ARN", policyDetail.ARN),
		)).Find(&policyDetail.Groups)
		rDB.Select(q.And(
			q.Eq("Type", "Users"),
			q.Eq("ARN", policyDetail.ARN),
		)).Find(&policyDetail.Users)
		rDB.Select(q.And(
			q.Eq("Type", "Roles"),
			q.Eq("ARN", policyDetail.ARN),
		)).Find(&policyDetail.Roles)
		policyDetails = append(policyDetails, policyDetail)
	}
	return policyDetails
}

func GetAllPolicies(i *iam.IAM, region, inUse, inUseByRoles, inUseByUsers, inUseByGroups string, all bool) *models.GetResponseApi {
	var policies []models.IAMManagedPolicy
	var queries []q.Matcher
	if inUseByRoles != "" {
		queries = append(queries, q.Eq("InUseByRoles", inUseByRoles))
	}
	if inUseByUsers != "" {
		queries = append(queries, q.Eq("InUseByUsers", inUseByUsers))
	}
	if inUseByGroups != "" {
		queries = append(queries, q.Eq("InUseByGroups", inUseByGroups))
	}
	if region != "" {
		queries = append(queries, q.Eq("Region", region))
		if all {
			i.DB.Find("Region", region, &policies)
		} else if inUse != "" {
			i.DB.Select(q.And(
				q.Eq("InUse", inUse),
				q.Eq("Region", region),
			)).Find(&policies)
		} else {
			i.DB.Select(q.And(
				queries...,
			)).Find(&policies)
		}
	} else {
		if all {
			i.DB.All(&policies)
		} else if inUse != "" {
			i.DB.Find("InUse", inUse, &policies)
		} else {
			i.DB.Select(q.And(
				queries...,
			)).Find(&policies)
		}
	}

	policyDetails := getAllPolicyAttachments(i, policies)

	apiResponse := models.GetResponseApi{
		Data:  policyDetails,
		Count: len(policyDetails),
	}
	return &apiResponse
}

func getAllInstanceProfileResources(i *iam.IAM, profiles []models.IAMInstanceProfile) []models.IAMInstanceProfileDetails {
	var profileDetails []models.IAMInstanceProfileDetails
	for _, profile := range profiles {
		profileDetail := models.IAMInstanceProfileDetails{
			ARN:              profile.ARN,
			ID:               profile.ID,
			InUse:            profile.InUse,
			InUseByRoles:     profile.InUseByRoles,
			InUseByInstances: profile.InUseByInstances,
			Name:             profile.Name,
			Region:           profile.Region,
		}
		rDB := i.DB.From("IAMInstanceProfile")
		rDB.Select(q.And(
			q.Eq("Type", shared.Role),
			q.Eq("ARN", profileDetail.ARN),
		)).Find(&profileDetail.Roles)
		rDB.Select(q.And(
			q.Eq("Type", shared.Instance),
			q.Eq("ARN", profileDetail.ARN),
		)).Find(&profileDetail.Instances)
		rDB.Select(q.And(
			q.Eq("Type", shared.LaunchConfiguration),
			q.Eq("ARN", profileDetail.ARN),
		)).Find(&profileDetail.LaunchConfigurations)
		profileDetails = append(profileDetails, profileDetail)
	}
	return profileDetails
}

func GetAllInstanceProfiles(i *iam.IAM, region, inUse, inUseByRoles, inUseByInstances, inUseByLCs string, all bool) *models.GetResponseApi {
	var profiles []models.IAMInstanceProfile
	var queries []q.Matcher
	if inUseByRoles != "" {
		queries = append(queries, q.Eq("InUseByRoles", inUseByRoles))
	}
	if inUseByInstances != "" {
		queries = append(queries, q.Eq("InUseByInstances", inUseByInstances))
	}
	if inUseByLCs != "" {
		queries = append(queries, q.Eq("InUseByLCs", inUseByLCs))
	}
	if region != "" {
		queries = append(queries, q.Eq("Region", region))
		if all {
			i.DB.Find("Region", region, &profiles)
		} else if inUse != "" {
			i.DB.Select(q.And(
				q.Eq("InUse", inUse),
				q.Eq("Region", region),
			)).Find(&profiles)
		} else {
			i.DB.Select(q.And(
				queries...,
			)).Find(&profiles)
		}
	} else {
		if all {
			i.DB.All(&profiles)
		} else if inUse != "" {
			i.DB.Find("InUse", inUse, &profiles)
		} else {
			i.DB.Select(q.And(
				queries...,
			)).Find(&profiles)
		}
	}

	profileDetails := getAllInstanceProfileResources(i, profiles)

	apiResponse := models.GetResponseApi{
		Data:  profileDetails,
		Count: len(profileDetails),
	}
	return &apiResponse
}
