package iam

import (
	"sync"
	"testing"

	"github.com/asdine/storm/q"
	"github.com/aws/aws-sdk-go/aws"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func init() {
	shared.DBC = shared.PrepareDb("test.db")
	shared.DefaultRegion = "us-west-2"
}

func TestUsers(t *testing.T) {
	encodedPolicy := `{"Version": "2012-10-17","Statement": [{"Effect": "Allow","Action": "*","Resource": "*"}]}}]}`
	userArn := "arn:aws:iam::123456789:user/tester"
	userName := "tester"
	userId := "123456789"
	policyArn := "arn:aws:iam::123456789:policy/awsapp"
	policyName := "awsapp"
	iamManagedPolicies := []*awsiam.AttachedPolicy{
		&awsiam.AttachedPolicy{
			PolicyArn:  &policyArn,
			PolicyName: &policyName,
		},
	}
	inlinePolicyName := "Admin"
	decodedPolicy := DecodePolicy(encodedPolicy)
	iamInlinePolcies := []*awsiam.PolicyDetail{
		&awsiam.PolicyDetail{
			PolicyDocument: &decodedPolicy,
			PolicyName:     &inlinePolicyName,
		},
	}
	groupList := []*string{
		aws.String("Foo"),
		aws.String("Bar"),
		aws.String("Baz"),
	}
	userDetails := []*awsiam.UserDetail{
		&awsiam.UserDetail{
			Arn:                     aws.String(userArn),
			UserId:                  aws.String(userId),
			UserName:                aws.String(userName),
			GroupList:               groupList,
			UserPolicyList:          iamInlinePolcies,
			AttachedManagedPolicies: iamManagedPolicies,
		},
	}
	var wg sync.WaitGroup
	myiam := iam.DB(shared.DBC)
	wg.Add(1)
	go myiam.StoreUsers(userDetails, &wg)
	wg.Wait()
	var iamUser models.IAMUser
	var managedPolicyResource []models.IAMResource
	var inlinePolicyResource []models.IAMInlineResource

	resourceBucket := myiam.DB.From("IAMUser")

	//Get Policy
	myiam.DB.One("ARN", userArn, &iamUser)

	//Get Managed Policies for User
	resourceBucket.Select(q.And(
		q.Eq("ARN", userArn),
		q.Eq("ResourceID", policyArn),
	)).Find(&managedPolicyResource)

	//Get Inline Policies for User
	resourceBucket.Select(q.And(
		q.Eq("ARN", userArn),
		q.Eq("Name", inlinePolicyName),
	)).Find(&inlinePolicyResource)

	// Validate that the user was indeed stored into the database
	if (models.IAMUser{}) != iamUser {
		// Validate that the user matches against the ARN
		if iamUser.ARN != userArn {
			t.Errorf("ARN did not match %s", policyArn)
		}
		// Validate that the managed policy resource was indeed stored into the database
		if len(managedPolicyResource) > 0 {
			if managedPolicyResource[0].ResourceID != policyArn {
				t.Errorf("Failed to match against group Managed Policy id %s", policyArn)
			}
		} else {
			t.Error("Failed to store Managed Policy resource")
		}
		// Validate that the inline policy resource was indeed stored into the database
		if len(inlinePolicyResource) > 0 {
			if inlinePolicyResource[0].ARN != userArn && inlinePolicyResource[0].Name != inlinePolicyName {
				t.Error("Failed to match against Inline Policy resource %s", inlinePolicyName)
			}
		} else {
			t.Error("Failed to store Inline Policy resource")
		}
	} else {
		t.Errorf("User %s was not found", userName)
	}
}
