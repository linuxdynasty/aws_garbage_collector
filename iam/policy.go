package iam

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func (i *IAM) storeRolePolicyResources(policyArn string, resourceType string, resources []*iam.PolicyRole) {
	for _, resource := range resources {
		policyResource := models.IAMResource{
			ResourceID: *resource.RoleId,
			Name:       *resource.RoleName,
			Type:       resourceType,
			ARN:        policyArn,
		}
		resourceDB := i.DB.From("IAMManagedPolicy")
		if err := resourceDB.Save(&policyResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeUserPolicyResources(policyArn string, resourceType string, resources []*iam.PolicyUser) {
	for _, resource := range resources {
		policyResource := models.IAMResource{
			ResourceID: *resource.UserId,
			Name:       *resource.UserName,
			Type:       resourceType,
			ARN:        policyArn,
		}
		resourceDB := i.DB.From("IAMManagedPolicy")
		if err := resourceDB.Save(&policyResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeGroupPolicyResources(policyArn string, resourceType string, resources []*iam.PolicyGroup) {
	for _, resource := range resources {
		policyResource := models.IAMResource{
			ResourceID: *resource.GroupId,
			Name:       *resource.GroupName,
			Type:       resourceType,
			ARN:        policyArn,
		}
		resourceDB := i.DB.From("IAMManagedPolicy")
		if err := resourceDB.Save(&policyResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) fetchPolicyResources(policy *models.IAMManagedPolicy) error {
	params := &iam.ListEntitiesForPolicyInput{
		PolicyArn: &policy.ARN,
	}
	policy.InUseByGroups = "false"
	policy.InUseByUsers = "false"
	policy.InUseByRoles = "false"
	err := i.Client.ListEntitiesForPolicyPages(params,
		func(resp *iam.ListEntitiesForPolicyOutput, lastPage bool) bool {
			if len(resp.PolicyGroups) > 0 {
				i.storeGroupPolicyResources(policy.ARN, "Groups", resp.PolicyGroups)
				policy.InUseByGroups = "true"
			}
			if len(resp.PolicyUsers) > 0 {
				i.storeUserPolicyResources(policy.ARN, "Users", resp.PolicyUsers)
				policy.InUseByUsers = "true"
			}
			if len(resp.PolicyRoles) > 0 {
				i.storeRolePolicyResources(policy.ARN, "Roles", resp.PolicyRoles)
				policy.InUseByRoles = "true"
			}
			return true
		},
	)
	return err
}

func (i *IAM) storePolicies(policies []*iam.ManagedPolicyDetail, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, p := range policies {
		inUse := "false"
		if *p.AttachmentCount > 0 {
			inUse = "true"
		}
		policy := models.IAMManagedPolicy{
			ARN:              *p.Arn,
			ID:               *p.PolicyId,
			Name:             *p.PolicyName,
			VersionId:        *p.DefaultVersionId,
			AttachementCount: *p.AttachmentCount,
			Region:           i.Region,
			InUse:            inUse,
		}

		params := &iam.GetPolicyVersionInput{
			PolicyArn: &policy.ARN,
			VersionId: &policy.VersionId,
		}
		resp, err := i.Client.GetPolicyVersion(params)
		if err != nil {
			log.Fatal(err)
		}
		decodedData, _ := url.QueryUnescape(*resp.PolicyVersion.Document)
		var policyJson models.Policy
		json.Unmarshal([]byte(decodedData), &policyJson)
		policy.Policy = policyJson
		i.fetchPolicyResources(&policy)
		if err := i.DB.Save(&policy); err != nil {
			log.Fatal(err)
		}
	}
}
