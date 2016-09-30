package iam

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (i *IAM) fetchAndStoreInstanceProfiles(client *iam.IAM, wg *sync.WaitGroup) {
	defer wg.Done()
	params := &iam.ListInstanceProfilesInput{}
	err := client.ListInstanceProfilesPages(params,
		func(resp *iam.ListInstanceProfilesOutput, lastPage bool) bool {
			for _, ip := range resp.InstanceProfiles {
				iprofile := models.IAMInstanceProfile{
					ARN:              *ip.Arn,
					ID:               *ip.InstanceProfileId,
					Name:             *ip.InstanceProfileName,
					InUse:            "false",
					InUseByRoles:     "false",
					InUseByInstances: "false",
					Region:           *client.Config.Region,
				}
				if ip.Roles != nil {
					rDb := i.DB.From("IAMInstanceProfile")
					for _, role := range ip.Roles {
						ipResource := models.IAMResource{
							ResourceID: *role.RoleId,
							Name:       *role.RoleName,
							Type:       shared.Role,
							ARN:        iprofile.ARN,
						}
						if err := rDb.Save(&ipResource); err != nil {
							log.Fatal(err)
						}
					}
					if len(ip.Roles) > 0 {
						iprofile.InUseByRoles = "true"
						iprofile.InUse = "true"
					}
				}
				if err := i.DB.Save(&iprofile); err != nil {
					log.Fatal(err)
				}
			}
			return true
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
