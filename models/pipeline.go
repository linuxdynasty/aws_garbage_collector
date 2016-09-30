package models

type PipeLine struct {
	ID      string `storm:"id""`
	Name    string `storm:"index"`
	ImageId string `storm:"index"`
}
