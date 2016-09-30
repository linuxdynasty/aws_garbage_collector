package models

type Policy interface{}

type IAMInlineResource struct {
	ID       int    `storm:"id" json:"-"`
	Document Policy `json:"document"`
	Name     string `json:"name"`
	Type     string `storm:"index" json:"type"`
	ARN      string `storm:"index" json:"policy_id"`
}

type IAMResource struct {
	ID         int    `storm:"id" json:"-"`
	ResourceID string `storm:"index" json:"resource_id"`
	Name       string `json:"name"`
	Type       string `storm:"index" json:"type"`
	ARN        string `storm:"index" json:"policy_id"`
}

type IAMManagedPolicy struct {
	ARN              string `storm:"index" json:"arn"`
	AttachementCount int64  `storm:"index" json:"attachment_count"`
	ID               string `storm:"id" json:"id"`
	InUse            string `storm:"index" json:"in_use"`
	InUseByUsers     string `storm:"index" json:"in_use_by_users"`
	InUseByRoles     string `storm:"index" json:"in_use_by_roles"`
	InUseByGroups    string `storm:"index" json:"in_use_by_groups"`
	Name             string `storm:"index" json:"name"`
	VersionId        string `json:"version_id"`
	Policy           Policy `json:"policy"`
	Region           string `storm:"index" json:"region"`
}

type IAMInstanceProfile struct {
	ARN              string `storm:"index" json:"arn"`
	ID               string `storm:"id" json:"id"`
	InUse            string `storm:"index" json:"in_use"`
	InUseByRoles     string `storm:"index" json:"in_use_by_roles"`
	InUseByInstances string `storm:"index" json:"in_use_by_instances"`
	InUseByLCs       string `storm:"index" json:"in_use_by_launch_configurations"`
	Name             string `storm:"index" json:"name"`
	Region           string `storm:"index" json:"region"`
}

type IAMRole struct {
	ARN                  string `storm:"index" json:"arn"`
	AssumeRoleDocument   Policy `json:"assume_role_document"`
	ID                   string `storm:"id" json:"id"`
	InUse                string `storm:"index" json:"in_use"`
	ManagedPolicyCount   int64  `storm:"index" json:"managed_policy_count"`
	InlinePolicyCount    int64  `storm:"index" json:"inline_policy_count"`
	InstanceProfileCount int64  `storm:"index" json:"instance_profile_count"`
	Name                 string `storm:"index" json:"name"`
	Region               string `storm:"index" json:"region"`
}

type IAMUser struct {
	ARN                string `storm:"index" json:"arn"`
	ID                 string `storm:"id" json:"id"`
	GroupsCount        int64  `storm:"index" json:"groups_count"`
	ManagedPolicyCount int64  `storm:"index" json:"managed_policy_count"`
	InlinePolicyCount  int64  `storm:"index" json:"inline_policy_count"`
	InUse              string `storm:"index" json:"in_use"`
	Name               string `storm:"index" json:"name"`
	Region             string `storm:"index" json:"region"`
}
