package redshift

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (r *RedShift) redShiftCluster(client *redshift.Redshift) error {
	rDb := r.DB.From("SecurityGroup")
	params := &redshift.DescribeClustersInput{}
	err := client.DescribeClustersPages(params,
		func(resp *redshift.DescribeClustersOutput, lastPage bool) bool {
			for _, val := range resp.Clusters {
				for _, ec2Sg := range val.VpcSecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   *val.ClusterIdentifier,
						Name:         *val.ClusterIdentifier,
						ResourceType: shared.Cluster,
						AWSType:      shared.RedShift,
						GroupId:      *ec2Sg.VpcSecurityGroupId,
					}
					if err := rDb.Save(&resource); err != nil {
						log.Fatal(err)
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

func (r *RedShift) redShift(client *redshift.Redshift) error {
	rDb := r.DB.From("SecurityGroup")
	params := &redshift.DescribeClusterSecurityGroupsInput{}
	err := client.DescribeClusterSecurityGroupsPages(params,
		func(resp *redshift.DescribeClusterSecurityGroupsOutput, lastPage bool) bool {
			for _, val := range resp.ClusterSecurityGroups {
				for _, ec2Sg := range val.EC2SecurityGroups {
					group := &models.SecurityGroup{}
					if err := r.DB.One("Name", *ec2Sg.EC2SecurityGroupName, group); err != nil {
						log.Print(*ec2Sg.EC2SecurityGroupName, " - not found")
					} else {
						resource := models.SecurityGroupResource{
							ResourceId:   *val.ClusterSecurityGroupName,
							Name:         *val.ClusterSecurityGroupName,
							ResourceType: shared.ResourceSecurityGroup,
							AWSType:      shared.RedShift,
							GroupId:      group.ID,
						}
						if err := rDb.Save(&resource); err != nil {
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

func (r *RedShift) RedShift(region string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	// Create AWS RedShift Session
	session := session.New(&aws.Config{Region: &region})
	client := redshift.New(session)

	err = r.redShift(client)
	if err != nil {
		return err
	}
	r.redShiftCluster(client)
	if err != nil {
		return err
	}
	return nil
}
