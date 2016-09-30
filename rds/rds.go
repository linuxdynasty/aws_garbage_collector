package rds

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (r *RDS) rdsCluster(client *rds.RDS) error {
	rDb := r.DB.From("SecurityGroup")
	params := &rds.DescribeDBClustersInput{}
	resp, err := client.DescribeDBClusters(params)
	for _, val := range resp.DBClusters {
		for _, ec2Sg := range val.VpcSecurityGroups {
			resource := models.SecurityGroupResource{
				ResourceId:   *val.DBClusterArn,
				Name:         *val.DBClusterIdentifier,
				ResourceType: shared.Cluster,
				AWSType:      shared.Rds,
				GroupId:      *ec2Sg.VpcSecurityGroupId,
			}
			if err := rDb.Save(&resource); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *RDS) rdsInstance(client *rds.RDS) error {
	rDb := r.DB.From("SecurityGroup")
	params := &rds.DescribeDBInstancesInput{}
	err := client.DescribeDBInstancesPages(params,
		func(resp *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			for _, val := range resp.DBInstances {
				for _, ec2Sg := range val.VpcSecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   *val.DBInstanceArn,
						Name:         *val.DBInstanceIdentifier,
						ResourceType: shared.Instance,
						AWSType:      shared.Rds,
						GroupId:      *ec2Sg.VpcSecurityGroupId,
					}
					resource.GroupId = *ec2Sg.VpcSecurityGroupId
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

func (r *RDS) rds(client *rds.RDS) error {
	rDb := r.DB.From("SecurityGroup")
	params := &rds.DescribeDBSecurityGroupsInput{}
	err := client.DescribeDBSecurityGroupsPages(params,
		func(resp *rds.DescribeDBSecurityGroupsOutput, lastPage bool) bool {
			for _, val := range resp.DBSecurityGroups {
				for _, ec2Sg := range val.EC2SecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   *val.DBSecurityGroupArn,
						Name:         *val.DBSecurityGroupName,
						ResourceType: shared.ResourceSecurityGroup,
						AWSType:      shared.ResourceSecurityGroup,
						GroupId:      *ec2Sg.EC2SecurityGroupId,
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

func (r *RDS) RDS(region string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	// Create AWS RDS Session
	session := session.New(&aws.Config{Region: &region})
	client := rds.New(session)

	err = r.rds(client)
	if err != nil {
		return err
	}
	r.rdsCluster(client)
	if err != nil {
		return err
	}
	r.rdsInstance(client)
	if err != nil {
		return err
	}
	return nil
}
