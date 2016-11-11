package iam

import (
	"net/url"
	"sync"
	"testing"

	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func init() {
	shared.DBC = shared.PrepareDb("test.db")
	shared.DefaultRegion = "us-west-2"
}

func DecodePolicy(policy string) (decodedPolicy string) {
	decodedPolicy = url.QueryEscape(policy)
	return decodedPolicy
}

func TestRoles(t *testing.T) {
	encodedPolicy := `{"Version": "2012-10-17","Statement": [{"Effect": "Allow","Action": "*","Resource": "*"}]}}]}`
	//decodedPolicy := DecodePolicy(encodedPolicy)
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
	instanceProfileArn := "arn:aws:iam::123456789:instance-profile/aws-app-development"
	instanceProfileName := "aws-app-development"
	iamInstanceProfiles := []*awsiam.InstanceProfile{
		&awsiam.InstanceProfile{
			Arn:                 &instanceProfileArn,
			InstanceProfileName: &instanceProfileName,
		},
	}
	roleArn := "arn:aws:iam::123456789:role/aws-app-development"
	roleName := "aws-app-development"
	roleId := "foobar"
	rolePolicyDocument := "derf"
	roleDetails := []*awsiam.RoleDetail{
		&awsiam.RoleDetail{
			Arn: &roleArn,
			AssumeRolePolicyDocument: &rolePolicyDocument,
			AttachedManagedPolicies:  iamManagedPolicies,
			RolePolicyList:           iamInlinePolcies,
			InstanceProfileList:      iamInstanceProfiles,
			RoleName:                 &roleName,
			RoleId:                   &roleId,
		},
	}
	var wg sync.WaitGroup
	myiam := iam.DB(shared.DBC)
	wg.Add(1)
	myiam.StoreRoles(roleDetails, &wg)
	wg.Wait()
	var iamRole models.IAMRole
	myiam.DB.One("ARN", roleArn, &iamRole)
	if iamRole != (models.IAMRole{}) {
		if iamRole.ARN != roleArn {
			t.Errorf("ARN did not match %s", roleArn)
		}
	} else {
		t.Errorf("Failed to store role %s", iamRole.Name)
	}
	var iamManagedPolicyResource models.IAMResource
	var iamInlinePolicyResource models.IAMInlineResource
	resourceBucket := myiam.DB.From("IAMRole")
	resourceBucket.One("ResourceID", policyArn, iamManagedPolicyResource)
	if (models.IAMResource{}) != iamManagedPolicyResource {
		if iamManagedPolicyResource.ResourceID != policyArn {
			t.Errorf("Failed to store resource Managed Policy %s", policyArn)
		}
		if iamManagedPolicyResource.ARN != roleArn {
			t.Errorf("Role ARN did not match %s inside Managed Policy", roleArn)
		}
	}
	resourceBucket.One("Name", inlinePolicyName, iamInlinePolicyResource)
	if (models.IAMInlineResource{}) != iamInlinePolicyResource {
		if iamInlinePolicyResource.Name != inlinePolicyName {
			t.Errorf("Failed to store resource Inline Policy %s", inlinePolicyName)
		}
		if iamInlinePolicyResource.ARN != roleArn {
			t.Errorf("Role ARN did not match %s inside Inline Policy", roleArn)
		}
	}
}
