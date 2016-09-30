package shared

import "github.com/asdine/storm"

var (
	DBC                   *storm.DB
	DefaultRegion         string
	Cluster               string = "Cluster"
	DataPipeline          string = "Data Pipeline"
	Instance              string = "Instance"
	SecGroup              string = "Security Group"
	ResourceSecurityGroup string = "Resource Security Group"
	Ec2                   string = "EC2"
	ElastiCache           string = "Elasticache"
	Elb                   string = "ELB"
	Elbv2                 string = "ELBV2"
	LoadBalancer          string = "Load Balancer"
	RedShift              string = "Redshift"
	Rds                   string = "RDS"
	LaunchConfiguration   string = "Launch Configuration"
	NetworkInterface      string = "Network Interface"
	Role                  string = "Role"
	Policy                string = "Policy"
	ManagedPolicy         string = "Managed Policy"
	InlinePolicy          string = "Inline Policy"
	InstanceProfile       string = "Instance Profile"
	User                  string = "User"
	Group                 string = "Group"
)
