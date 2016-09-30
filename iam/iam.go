package iam

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func (i *IAM) fetchAndStoreIAMResources(client *iam.IAM) error {
	params := &iam.GetAccountAuthorizationDetailsInput{
		Filter: []*string{
			aws.String("LocalManagedPolicy"),
		},
	}
	err := client.GetAccountAuthorizationDetailsPages(params,
		func(resp *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) bool {
			var wg sync.WaitGroup
			wg.Add(4)
			go i.storeRoles(client, resp.RoleDetailList, &wg)
			go i.storeUsers(client, resp.UserDetailList, &wg)
			go i.storePolicies(client, resp.Policies, &wg)
			go i.fetchAndStoreInstanceProfiles(client, &wg)
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
	svc := iam.New(session)
	if err := i.fetchAndStoreIAMResources(svc); err != nil {
		log.Fatal(err)
	}
	return err
}
