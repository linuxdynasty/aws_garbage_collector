package models

type SecurityGroupDetails struct {
	ID                   string                  `json:"id"`
	Name                 string                  `json:"name"`
	Description          string                  `json:"description"`
	InUse                string                  `json:"in_use"`
	InUseBySGOnly        string                  `json:"in_use_by_sg_only"`
	Region               string                  `json:"region"`
	Tags                 []Tag                   `json:"tags"`
	EC2                  []SecurityGroupResource `json:"ec2"`
	ElastiCache          []SecurityGroupResource `json:"elasticache"`
	ELB                  []SecurityGroupResource `json:"elb"`
	ELBV2                []SecurityGroupResource `json:"elbv2"`
	LaunchConfiguration  []SecurityGroupResource `json:"launch_configuration"`
	RDS                  []SecurityGroupResource `json:"rds"`
	RedShift             []SecurityGroupResource `json:"redshift"`
	SourceSecurityGroups []SourceSecurityGroup   `json:"source_security_groups"`
}

type LaunchConfigDetails struct {
	ID              string           `storm:"id" json:"id"`
	Name            string           `json:"name"`
	InUse           string           `storm:"index" json:"in_use"`
	ImageId         string           `storm:"index" json:"image_id"`
	InstanceProfile string           `storm:"index" json:"instance_profile"`
	Region          string           `storm:"index" json:"region"`
	AutoScaleGroups []AutoScaleGroup `json:"auto_scale_groups"`
}

type EC2AmiDetails struct {
	ID                  string                `json:"id"`
	Name                string                `json:"name"`
	Description         string                `json:"description"`
	ImageLocation       string                `json:"image_location"`
	Region              string                `json:"region"`
	VirtualizationType  string                `json:"virt_type"`
	State               string                `json:"state"`
	SnapshotIds         []string              `json:"snapshot_ids"`
	InUse               string                `json:"in_use"`
	InUseByLC           string                `json:"in_use_by_lc"`
	InUseByInstance     string                `json:"in_use_by_instance"`
	InUseByDataPipeline string                `json:"in_use_by_data_pipeline"`
	DataPipelines       []PipeLine            `json:"data_pipelines"`
	EC2                 []EC2Instance         `json:"instances"`
	LC                  []LaunchConfiguration `json:"launch_configurations"`
	Tags                []Tag
}

type IAMManagedPolicyDetails struct {
	ARN              string        `storm:"index" json:"arn"`
	AttachementCount int64         `storm:"index" json:"attachment_count"`
	ID               string        `storm:"id" json:"id"`
	InUse            string        `storm:"index" json:"in_use"`
	InUseByUsers     string        `storm:"index" json:"in_use_by_users"`
	InUseByRoles     string        `storm:"index" json:"in_use_by_roles"`
	InUseByGroups    string        `storm:"index" json:"in_use_by_groups"`
	Name             string        `storm:"index" json:"name"`
	VersionId        string        `json:"version_id"`
	Policy           Policy        `json:"policy"`
	Region           string        `storm:"index" json:"region"`
	Groups           []IAMResource `json:"groups"`
	Roles            []IAMResource `json:"roles"`
	Users            []IAMResource `json:"users"`
}

type IAMInstanceProfileDetails struct {
	ARN                  string        `storm:"index" json:"arn"`
	ID                   string        `storm:"id" json:"id"`
	InUse                string        `storm:"index" json:"in_use"`
	InUseByRoles         string        `storm:"index" json:"in_use_by_roles"`
	InUseByInstances     string        `storm:"index" json:"in_use_by_instances"`
	Name                 string        `storm:"index" json:"name"`
	Region               string        `storm:"index" json:"region"`
	Roles                []IAMResource `json:"roles"`
	Instances            []IAMResource `json:"instances"`
	LaunchConfigurations []IAMResource `json:"launch_configurations"`
}
