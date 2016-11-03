package iam

import (
	"testing"

	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func init() {
	shared.DBC = shared.PrepareDb("test.db")
	shared.DefaultRegion = "us-west-2"
}

func TestReader(t *testing.T) {

}
