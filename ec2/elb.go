package ec2

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func (e *EC2) storeELBResource(client *elb.ELB) error {
	params := &elb.DescribeLoadBalancersInput{}
	rDb := e.DB.From("SecurityGroup")
	err := client.DescribeLoadBalancersPages(params,
		func(resp *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, elb := range resp.LoadBalancerDescriptions {
				instanceName := *elb.LoadBalancerName
				for _, val := range elb.SecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   instanceName,
						Name:         instanceName,
						ResourceType: shared.LoadBalancer,
						AWSType:      shared.Elb,
						GroupId:      *val,
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

func (e *EC2) storeELBv2Resource(client *elbv2.ELBV2) error {
	params := &elbv2.DescribeLoadBalancersInput{}
	rDb := e.DB.From("SecurityGroup")
	err := client.DescribeLoadBalancersPages(params,
		func(resp *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, elbv2 := range resp.LoadBalancers {
				instanceName := *elbv2.LoadBalancerName
				for _, val := range elbv2.SecurityGroups {
					resource := models.SecurityGroupResource{
						ResourceId:   *elbv2.LoadBalancerArn,
						Name:         instanceName,
						ResourceType: shared.LoadBalancer,
						AWSType:      shared.Elbv2,
						GroupId:      *val,
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

func (e *EC2) StoreELBs(region string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Create AWS ELB and ELBv2 Sessions
	session := session.New(&aws.Config{Region: &region})
	elbClient := elb.New(session)
	elbv2Client := elbv2.New(session)

	if err := e.storeELBResource(elbClient); err != nil {
		log.Fatal(err)
	}
	if err := e.storeELBv2Resource(elbv2Client); err != nil {
		log.Fatal(err)
	}
}
