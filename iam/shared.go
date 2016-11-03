package iam

import (
	"github.com/asdine/storm"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMReader interface {
	GetAccountAuthorizationDetailsPages(input *iam.GetAccountAuthorizationDetailsInput, fn func(output *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool)) error
	ListInstanceProfilesPages(input *iam.ListInstanceProfilesInput, fn func(output *iam.ListInstanceProfilesOutput, lastPage bool) (shouldContinue bool)) error
	ListEntitiesForPolicyPages(input *iam.ListEntitiesForPolicyInput, fn func(output *iam.ListEntitiesForPolicyOutput, lastPage bool) (shouldContinue bool)) error
	GetPolicyVersion(input *iam.GetPolicyVersionInput) (*iam.GetPolicyVersionOutput, error)
}

type IAM struct {
	DB     *storm.DB
	Client IAMReader
	Region string
}

func DB(db *storm.DB) IAM {
	iam := IAM{
		DB: db,
	}
	return iam
}
