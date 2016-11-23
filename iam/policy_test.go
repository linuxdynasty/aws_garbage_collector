package iam

import (
	"fmt"
	"sync"
	"testing"

	"github.com/asdine/storm/q"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

type FakeIAMClient struct{}

func (f *FakeIAMClient) GetAccountAuthorizationDetailsPages(input *awsiam.GetAccountAuthorizationDetailsInput, fn func(output *awsiam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool)) error {
	return nil
}

func (f *FakeIAMClient) ListInstanceProfilesPages(input *awsiam.ListInstanceProfilesInput, fn func(output *awsiam.ListInstanceProfilesOutput, lastPage bool) (shouldContinue bool)) error {
	instanceProfileId := "123456789"
	instanceProfileArn := "arn:aws:iam::123456789:instance-profile/aws-app-development"
	instanceProfileName := "aws-app-development"
	roleId := "987654321"
	//roleArn := "arn:aws:iam::123456789:role/aws-app-development"
	roleName := "aws-app-development"
	roles := []*awsiam.Role{
		&awsiam.Role{
			//Arn:      &roleArn,
			RoleId:   &roleId,
			RoleName: &roleName,
		},
	}
	instanceProfiles := []*awsiam.InstanceProfile{
		&awsiam.InstanceProfile{
			InstanceProfileId:   &instanceProfileId,
			Arn:                 &instanceProfileArn,
			InstanceProfileName: &instanceProfileName,
			Roles:               roles,
		},
	}
	output := &awsiam.ListInstanceProfilesOutput{
		InstanceProfiles: instanceProfiles,
	}
	fmt.Println(output)

	fn(output, true)
	return nil
}

func (f *FakeIAMClient) ListEntitiesForPolicyPages(input *awsiam.ListEntitiesForPolicyInput, fn func(output *awsiam.ListEntitiesForPolicyOutput, lastPage bool) (shouldContinue bool)) error {
	//Group settings
	groupId := "987654321"
	groupName := "tester"
	policyGroups := []*awsiam.PolicyGroup{
		&awsiam.PolicyGroup{
			GroupId:   &groupId,
			GroupName: &groupName,
		},
	}

	//Role settings
	roleId := "654321789"
	roleName := "app-tester"
	policyRoles := []*awsiam.PolicyRole{
		&awsiam.PolicyRole{
			RoleId:   &roleId,
			RoleName: &roleName,
		},
	}

	//User settings
	userId := "12341234"
	userName := "appy"
	policyUsers := []*awsiam.PolicyUser{
		&awsiam.PolicyUser{
			UserId:   &userId,
			UserName: &userName,
		},
	}
	output := &awsiam.ListEntitiesForPolicyOutput{
		PolicyGroups: policyGroups,
		PolicyRoles:  policyRoles,
		PolicyUsers:  policyUsers,
	}
	fn(output, true)
	return nil
}

func (f *FakeIAMClient) GetPolicyVersion(input *awsiam.GetPolicyVersionInput) (*awsiam.GetPolicyVersionOutput, error) {
	encodedPolicy := `{"Version": "2012-10-17","Statement": [{"Effect": "Allow","Action": "*","Resource": "*"}]}`
	decodedPolicy := DecodePolicy(encodedPolicy)
	policyVersion := awsiam.PolicyVersion{
		Document: &decodedPolicy,
	}
	output := awsiam.GetPolicyVersionOutput{
		PolicyVersion: &policyVersion,
	}

	return &output, nil
}

func init() {
	shared.DBC = shared.PrepareDb("test.db")
	shared.DefaultRegion = "us-west-2"
}

func TestPolicies(t *testing.T) {
	//decodedPolicy := DecodePolicy(encodedPolicy)
	policyArn := "arn:aws:iam::123456789:policy/awsapp"
	policyName := "awsapp"
	policyId := "ADJKAHJD233"
	policyDefaultVersion := "v1"
	attachmentCount := int64(1)
	iamManagedPolicies := []*awsiam.ManagedPolicyDetail{
		&awsiam.ManagedPolicyDetail{
			Arn:              &policyArn,
			PolicyName:       &policyName,
			PolicyId:         &policyId,
			DefaultVersionId: &policyDefaultVersion,
			AttachmentCount:  &attachmentCount,
		},
	}
	var wg sync.WaitGroup
	myiam := iam.DB(shared.DBC)
	client := &FakeIAMClient{}
	myiam.Client = client
	myiam.Region = "us-west-2"
	wg.Add(1)
	go myiam.StorePolicies(iamManagedPolicies, &wg)
	wg.Wait()
	//myiam.FetchPolicyResources()
	var iamPolicy models.IAMManagedPolicy
	var groupPolicyResource []models.IAMResource
	var userPolicyResource []models.IAMResource
	var rolePolicyResource []models.IAMResource

	resourceBucket := myiam.DB.From("IAMManagedPolicy")

	//Get Policy
	myiam.DB.One("ARN", policyArn, &iamPolicy)

	//Get Groups for Policy
	resourceBucket.Select(q.And(
		q.Eq("ARN", policyArn),
		q.Eq("Type", "Groups"),
	)).Find(&groupPolicyResource)

	//Get Users for Policy
	resourceBucket.Select(q.And(
		q.Eq("ARN", policyArn),
		q.Eq("Type", "Users"),
	)).Find(&userPolicyResource)

	//Get Roles for Policy
	resourceBucket.Select(q.And(
		q.Eq("ARN", policyArn),
		q.Eq("Type", "Roles"),
	)).Find(&rolePolicyResource)

	// Validate that the policy was indeed stored into the database
	if (models.IAMManagedPolicy{}) != iamPolicy {
		// Validate that the policy matches against the ARN
		if iamPolicy.ARN != policyArn {
			t.Errorf("ARN did not match %s", policyArn)
		}
		// Validate that the group resource was indeed stored into the database
		if len(groupPolicyResource) > 0 {
			//if (models.IAMResource{}) != groupPolicyResource {
			if groupPolicyResource[0].ResourceID != "987654321" {
				t.Error("Failed to match against group resource id 987654321")
			}
		} else {
			t.Error("Failed to store group resource")
		}
		// Validate that the user resource was indeed stored into the database
		if len(userPolicyResource) > 0 {
			//if (models.IAMResource{}) != userPolicyResource {
			if userPolicyResource[0].ResourceID != "12341234" {
				t.Error("Failed to match against user resource id 12341234")
			}
		} else {
			t.Error("Failed to store user resource")
		}
		// Validate that the role resource was indeed stored into the database
		if len(rolePolicyResource) > 0 {
			//if (models.IAMResource{}) != rolePolicyResource {
			if rolePolicyResource[0].ResourceID != "654321789" {
				t.Error("Failed to match against role resource id 654321789")
			}
		} else {
			t.Error("Failed to store role resource")
		}

	} else {
		t.Errorf("Failed to store policy %s", iamPolicy.Name)
	}
}
