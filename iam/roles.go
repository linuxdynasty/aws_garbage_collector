package iam

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (i *IAM) storeInstanceProfileRoleResources(roleArn string, resources []*iam.InstanceProfile) {
	for _, resource := range resources {
		RoleResource := models.IAMResource{
			ResourceID: *resource.Arn,
			Name:       *resource.InstanceProfileName,
			Type:       shared.InstanceProfile,
			ARN:        roleArn,
		}
		resourceDB := i.DB.From("IAMRole")
		if err := resourceDB.Save(&RoleResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeInlinePolicyRoleResources(roleArn string, resources []*iam.PolicyDetail) {
	for _, resource := range resources {
		RoleResource := models.IAMInlineResource{
			Document: *resource.PolicyDocument,
			Name:     *resource.PolicyName,
			Type:     shared.InlinePolicy,
			ARN:      roleArn,
		}
		resourceDB := i.DB.From("IAMRole")
		if err := resourceDB.Save(&RoleResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeManagedPolicyRoleResources(roleArn string, resources []*iam.AttachedPolicy) {
	for _, resource := range resources {
		RoleResource := models.IAMResource{
			ResourceID: *resource.PolicyArn,
			Name:       *resource.PolicyName,
			Type:       shared.ManagedPolicy,
			ARN:        roleArn,
		}
		resourceDB := i.DB.From("IAMRole")
		if err := resourceDB.Save(&RoleResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeRoles(roles []*iam.RoleDetail, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, r := range roles {
		role := models.IAMRole{
			ARN:                  *r.Arn,
			AssumeRoleDocument:   *r.AssumeRolePolicyDocument,
			ID:                   *r.RoleId,
			Name:                 *r.RoleName,
			Region:               i.Region,
			ManagedPolicyCount:   int64(len(r.AttachedManagedPolicies)),
			InlinePolicyCount:    int64(len(r.RolePolicyList)),
			InstanceProfileCount: int64(len(r.InstanceProfileList)),
		}
		role.InUse = "false"
		if role.ManagedPolicyCount > 0 || role.ManagedPolicyCount > 0 || role.InstanceProfileCount > 0 {
			role.InUse = "true"
		}
		if err := i.DB.Save(&role); err != nil {
			log.Fatal(err)
		}
		i.storeManagedPolicyRoleResources(role.ARN, r.AttachedManagedPolicies)
		i.storeInlinePolicyRoleResources(role.ARN, r.RolePolicyList)
		i.storeInstanceProfileRoleResources(role.ARN, r.InstanceProfileList)
	}
}
