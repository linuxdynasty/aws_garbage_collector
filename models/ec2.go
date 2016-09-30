package models

type EC2Instance struct {
	ID              string `storm:"id""`
	Name            string `storm:"index"`
	Description     string
	ImageId         string `storm:"index"`
	InstanceProfile string `storm:"index"`
	Region          string `storm:"index"`
	SSHKey          string
	State           string
}

type EC2Ami struct {
	Description         string
	ID                  string `storm:"id"`
	ImageLocation       string
	InUse               string
	InUseByInstance     string
	InUseByDataPipeline string
	InUseByLC           string
	Name                string `storm:"index"`
	Region              string `storm:"index"`
	SnapshotIds         []string
	State               string
	VirtualizationType  string `storm:"index"`
}

type EC2NetworkInterface struct {
	ID          string `storm:"id"`
	Name        string `storm:"index"`
	Description string
	Region      string `storm:"index"`
	Status      string
	PrivateIp   string
}

type SecurityGroupResource struct {
	ID           int    `storm:"id" json:"-"`
	AWSType      string `storm:"index" json:"aws_type"`
	GroupId      string `storm:"index" json:"-"`
	Name         string `json:"name"`
	ResourceId   string `storm:"index" json:"resource_id,omitempty"`
	ResourceType string `storm:"index" json:"resource_type"`
}

type SourceSecurityGroup struct {
	ID            int    `storm:"id" json:"-"`
	GroupId       string `storm:"index" json:"-"`
	SourceGroupId string `storm:"index" json:"source_group_id"`
	Name          string `json:"name"`
}

type SecurityGroup struct {
	ID            string `storm:"id" json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	InUse         string `storm:"index" json:"in_use"`
	InUseBySGOnly string `storm:"index" json:"in_use_by_sg_only"`
	Region        string `storm:"index" json:"region"`
}
