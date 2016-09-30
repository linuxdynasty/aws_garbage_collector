package elasticache

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (e *ElastiCache) elastiCacheCluster(client *elasticache.ElastiCache) error {
	rDb := e.DB.From("SecurityGroup")
	params := &elasticache.DescribeCacheClustersInput{}
	err := client.DescribeCacheClustersPages(params,
		func(resp *elasticache.DescribeCacheClustersOutput, lastPage bool) bool {
			for _, val := range resp.CacheClusters {
				for _, ec2Sg := range val.SecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   *val.CacheClusterId,
						Name:         *val.CacheClusterId,
						ResourceType: shared.Cluster,
						AWSType:      shared.ElastiCache,
						GroupId:      *ec2Sg.SecurityGroupId,
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

func (e *ElastiCache) elastiCache(client *elasticache.ElastiCache) error {
	rDb := e.DB.From("SecurityGroup")
	params := &elasticache.DescribeCacheSecurityGroupsInput{}
	err := client.DescribeCacheSecurityGroupsPages(params,
		func(resp *elasticache.DescribeCacheSecurityGroupsOutput, lastPage bool) bool {
			for _, val := range resp.CacheSecurityGroups {
				for _, ec2Sg := range val.EC2SecurityGroups {
					group := &models.SecurityGroup{}
					if err := e.DB.One("Name", *ec2Sg.EC2SecurityGroupName, group); err != nil {
						log.Print(*ec2Sg.EC2SecurityGroupName, " - not found")
					} else {
						resource := models.SecurityGroupResource{
							ResourceId:   *val.CacheSecurityGroupName,
							Name:         *val.CacheSecurityGroupName,
							ResourceType: shared.ResourceSecurityGroup,
							AWSType:      shared.ElastiCache,
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

func (e *ElastiCache) ElastiCache(region string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	// Create AWS ElastiCache Session
	session := session.New(&aws.Config{Region: &region})
	client := elasticache.New(session)

	err = e.elastiCache(client)
	if err != nil {
		return err
	}
	err = e.elastiCacheCluster(client)
	if err != nil {
		return err
	}
	return nil
}
