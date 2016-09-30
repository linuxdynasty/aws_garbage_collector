package asgs

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (l *LC) launchConfigurations(client *autoscaling.AutoScaling) error {
	iDb := l.DB.From("IAMInstanceProfile")
	params := &autoscaling.DescribeLaunchConfigurationsInput{}
	err := client.DescribeLaunchConfigurationsPages(params,
		func(resp *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) bool {
			for _, launchconfiguration := range resp.LaunchConfigurations {
				lc := models.LaunchConfiguration{
					ID:      *launchconfiguration.LaunchConfigurationARN,
					Name:    *launchconfiguration.LaunchConfigurationName,
					ImageId: *launchconfiguration.ImageId,
					Region:  *client.Config.Region,
				}
				if launchconfiguration.IamInstanceProfile != nil {
					lc.InstanceProfile = *launchconfiguration.IamInstanceProfile
					ipResource := models.IAMResource{
						ResourceID: lc.ID,
						Name:       lc.Name,
						Type:       shared.LaunchConfiguration,
						ARN:        lc.InstanceProfile,
					}
					if err := iDb.Save(&ipResource); err != nil {
						log.Fatal(err)
					}
				}
				var asg models.AutoScaleGroup
				if findErr := l.ASGBucket.One("LaunchConfigurationName", lc.Name, &asg); findErr == nil {
					lc.InUse = "true"
				} else {
					lc.InUse = "false"
				}
				rDb := l.DB.From("SecurityGroups")
				for _, sg := range launchconfiguration.SecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   lc.ID,
						Name:         lc.Name,
						ResourceType: shared.LaunchConfiguration,
						AWSType:      shared.LaunchConfiguration,
						GroupId:      *sg,
					}
					if err := rDb.Save(&resource); err != nil {
						log.Fatal(err)
					}
				}
				if saveErr := l.DB.Save(&lc); saveErr != nil {
					log.Fatal(saveErr)
				}
			}
			return true
		},
	)
	return err
}

func (l *LC) asgs(client *autoscaling.AutoScaling) error {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	err := client.DescribeAutoScalingGroupsPages(params,
		func(resp *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			for _, autoscalegroups := range resp.AutoScalingGroups {
				asg := models.AutoScaleGroup{
					ID:   *autoscalegroups.AutoScalingGroupARN,
					Name: *autoscalegroups.AutoScalingGroupName,
					LaunchConfigurationName: *autoscalegroups.LaunchConfigurationName,
					Region:                  *client.Config.Region,
				}
				if err := l.ASGBucket.Save(&asg); err != nil {
					log.Fatal(err)
				}
			}
			return true
		})
	return err
}

func (l *LC) LaunchConfigurations(region string, wg *sync.WaitGroup) error {
	defer wg.Done()

	var err error
	session := session.New(&aws.Config{Region: &region})
	svc := autoscaling.New(session)

	err = l.asgs(svc)
	if err != nil {
		log.Fatal(err)
	}
	err = l.launchConfigurations(svc)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (l *LC) DeleteByName(lcNames []string, region string) []models.DeleteStatus {
	session := session.New(&aws.Config{Region: &region})
	svc := autoscaling.New(session)
	statuses := []models.DeleteStatus{}

	for _, lcName := range lcNames {
		var lc models.LaunchConfiguration
		params := &autoscaling.DeleteLaunchConfigurationInput{
			LaunchConfigurationName: aws.String(lcName),
		}
		_, err := svc.DeleteLaunchConfiguration(params)
		status := models.DeleteStatus{
			Name: lcName,
		}
		if queryErr := l.DB.One("Name", lcName, &lc); queryErr == nil {
			status.ID = lc.ID
		}
		if err != nil {
			status.Deleted = false
			status.Message = err.Error()
		} else {
			status.Deleted = true
			status.Message = fmt.Sprintf("%s deleted successfully", lcName)
			l.DB.DeleteStruct(&lc)
		}
		statuses = append(statuses, status)
	}
	return statuses
}
