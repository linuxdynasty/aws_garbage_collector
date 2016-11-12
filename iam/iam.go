package iam

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func (i *IAM) fetchAndStoreIAMResources() error {
	params := &iam.GetAccountAuthorizationDetailsInput{
		Filter: []*string{
			aws.String("LocalManagedPolicy"),
		},
	}
	err := i.Client.GetAccountAuthorizationDetailsPages(params,
		func(resp *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) bool {
			var wg sync.WaitGroup
			wg.Add(4)
			go i.StoreRoles(resp.RoleDetailList, &wg)
			go i.StoreUsers(resp.UserDetailList, &wg)
			go i.StorePolicies(resp.Policies, &wg)
			go i.fetchAndStoreInstanceProfiles(&wg)
			wg.Wait()

			return true
		},
	)
	return err
}

func (i *IAM) FetchAndStoreIAM(region string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error
	session := session.New(&aws.Config{Region: &region})
	i.Client = iam.New(session)
	i.Region = region
	if err := i.fetchAndStoreIAMResources(); err != nil {
		log.Fatal(err)
	}
	return err
}
