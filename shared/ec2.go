package shared

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func IsRegionValid(region, baseRegion string) bool {
	regions := Regions(baseRegion)
	var matchedRegion bool
	for _, r := range regions {
		if region == r {
			matchedRegion = true
			break
		}
	}
	return matchedRegion
}

func Regions(defaultRegion string) []string {
	var regions []string
	sess := session.New(&aws.Config{Region: &defaultRegion})
	svc := ec2.New(sess)
	params := &ec2.DescribeRegionsInput{}
	resp, err := svc.DescribeRegions(params)
	if err != nil {
		log.Fatal("Failed to collect regions", err.Error())
	}
	for _, region := range resp.Regions {
		regions = append(regions, *region.RegionName)
	}
	return regions
}
