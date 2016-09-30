package ec2

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (e *EC2) storeInstances(client *ec2.EC2) error {
	iDb := e.DB.From("IAMInstanceProfile")
	rBucket := e.DB.From("SecurityGroup")
	params := &ec2.DescribeInstancesInput{}
	err := client.DescribeInstancesPages(params,
		func(resp *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range resp.Reservations {
				for _, instance := range reservation.Instances {
					instanceName := *instance.PrivateDnsName
					for _, tag := range instance.Tags {
						if *tag.Key == "Name" {
							instanceName = *tag.Value
						}

						awstag := models.Tag{
							ResourceId: *instance.InstanceId,
							Key:        *tag.Key,
							Value:      *tag.Value,
							Region:     *client.Config.Region,
						}
						if err := e.DB.Save(&awstag); err != nil {
							log.Fatal(err)
						}
					}
					ec2Instance := models.EC2Instance{
						ID:      *instance.InstanceId,
						Name:    instanceName,
						ImageId: *instance.ImageId,
						Region:  *client.Config.Region,
						State:   *instance.State.Name,
					}
					if instance.IamInstanceProfile != nil {
						ec2Instance.InstanceProfile = *instance.IamInstanceProfile.Arn
						ipResource := models.IAMResource{
							ResourceID: ec2Instance.ID,
							Name:       ec2Instance.Name,
							Type:       shared.Instance,
							ARN:        ec2Instance.InstanceProfile,
						}
						if err := iDb.Save(&ipResource); err != nil {
							log.Fatal(err)
						}
					}
					if instance.KeyName != nil {
						ec2Instance.SSHKey = *instance.KeyName
					}
					if err := e.DB.Save(&ec2Instance); err != nil {
						log.Fatal(err)
					}
					for _, val := range instance.SecurityGroups {
						resource := models.SecurityGroupResource{
							ResourceId:   *instance.InstanceId,
							Name:         instanceName,
							ResourceType: "Instance",
							AWSType:      "EC2",
							GroupId:      *val.GroupId,
						}
						if err := rBucket.Save(&resource); err != nil {
							log.Fatal(err)
						}
					}
				}
			}
			return true
		})
	if err != nil {
		return err
	}
	return nil
}

func (e *EC2) storeNetworkInterfaces(client *ec2.EC2) error {
	params := &ec2.DescribeNetworkInterfacesInput{}
	resp, err := client.DescribeNetworkInterfaces(params)
	if err != nil {
		return err
	}
	for _, net := range resp.NetworkInterfaces {
		var instanceName string
		for _, tag := range net.TagSet {
			if *tag.Key == "Name" {
				instanceName = *tag.Value
			}

			awstag := models.Tag{
				ResourceId: *net.NetworkInterfaceId,
				Key:        *tag.Key,
				Value:      *tag.Value,
				Region:     *client.Config.Region,
			}
			if err := e.DB.Save(&awstag); err != nil {
				log.Fatal(err)
			}
		}
		netInterface := models.EC2NetworkInterface{
			ID:          *net.NetworkInterfaceId,
			Name:        instanceName,
			Description: *net.Description,
			Region:      *client.Config.Region,
			Status:      *net.Status,
		}
		if net.PrivateIpAddress != nil {
			netInterface.PrivateIp = *net.PrivateIpAddress
		}
		for _, val := range net.Groups {
			resource := models.SecurityGroupResource{
				ResourceId:   *net.NetworkInterfaceId,
				Name:         *net.Description,
				ResourceType: "Network Interface",
				AWSType:      "EC2",
			}
			resource.GroupId = *val.GroupId
			rBucket := e.DB.From("SecurityGroup")
			if err := rBucket.Save(&resource); err != nil {
				log.Fatal(err)
			}
		}
	}
	return err
}

func (e *EC2) FetchAndStoreEC2(region string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error
	session := session.New(&aws.Config{Region: &region})
	svc := ec2.New(session)

	if ec2Err := e.storeInstances(svc); ec2Err != nil {
		return err
	}
	if ec2Err := e.storeNetworkInterfaces(svc); ec2Err != nil {
		return err
	}
	return err
}
