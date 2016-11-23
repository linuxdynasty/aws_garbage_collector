package iam

import (
	"sync"
	"testing"

	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func init() {
	shared.DBC = shared.PrepareDb("test.db")
	shared.DefaultRegion = "us-west-2"
}

func TestInstanceProfiles(t *testing.T) {
	var wg sync.WaitGroup
	myiam := iam.DB(shared.DBC)
	client := &FakeIAMClient{}
	myiam.Client = client
	myiam.Region = "us-west-2"
	wg.Add(1)
	go myiam.StoreInstanceProfiles(&wg)
	wg.Wait()
}
