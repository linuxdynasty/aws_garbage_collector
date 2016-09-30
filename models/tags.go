package models

type Tag struct {
	ID         int    `storm:"id" json:"-"`
	ResourceId string `storm:"index" json:"-"`
	Key        string `storm:"index" json:"key"`
	Value      string `storm:"index" json:"value"`
	Region     string `storm:"index" json:"region"`
}
