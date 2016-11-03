package iam

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (i *IAM) storeUserManagedPolicyResources(userArn string, resources []*iam.AttachedPolicy) {
	for _, resource := range resources {
		UserResource := models.IAMResource{
			ResourceID: *resource.PolicyArn,
			Name:       *resource.PolicyName,
			Type:       shared.ManagedPolicy,
			ARN:        userArn,
		}
		resourceDB := i.DB.From("IAMUser")
		if err := resourceDB.Save(&UserResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeUserInlinePolicyResources(userArn string, resources []*iam.PolicyDetail) {
	for _, resource := range resources {
		UserResource := models.IAMInlineResource{
			Document: *resource.PolicyDocument,
			Name:     *resource.PolicyName,
			Type:     shared.InlinePolicy,
			ARN:      userArn,
		}
		resourceDB := i.DB.From("IAMUser")
		if err := resourceDB.Save(&UserResource); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *IAM) storeUsers(users []*iam.UserDetail, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, u := range users {
		user := models.IAMUser{
			ARN:                *u.Arn,
			ID:                 *u.UserId,
			Name:               *u.UserName,
			Region:             i.Region,
			GroupsCount:        int64(len(u.GroupList)),
			ManagedPolicyCount: int64(len(u.AttachedManagedPolicies)),
			InlinePolicyCount:  int64(len(u.UserPolicyList)),
		}
		user.InUse = "false"
		if user.GroupsCount > 0 || user.ManagedPolicyCount > 0 {
			user.InUse = "true"
		}
		if err := i.DB.Save(&user); err != nil {
			log.Fatal(err)
		}
		i.storeUserManagedPolicyResources(user.ARN, u.AttachedManagedPolicies)
		i.storeUserInlinePolicyResources(user.ARN, u.UserPolicyList)
	}
}
