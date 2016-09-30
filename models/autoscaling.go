package models

type LaunchConfiguration struct {
	ID              string `storm:"id" json:"id"`
	Name            string `json:"name"`
	InUse           string `storm:"index" json:"in_use"`
	ImageId         string `storm:"index" json:"image_id"`
	InstanceProfile string `storm:"index" json:"instance_profile"`
	Region          string `storm:"index" json:"region"`
}

type AutoScaleGroup struct {
	ID                      string `storm:"id" json:"id"`
	Name                    string `json:"name"`
	LaunchConfigurationName string `storm:"index" json:"launch_configuration_name"`
	Region                  string `storm:"index" json:"region"`
}
