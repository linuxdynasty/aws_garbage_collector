package datapipelines

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/datapipeline"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func (p *DataPipeline) listPipelines(client *datapipeline.DataPipeline) ([]models.PipeLine, error) {
	params := &datapipeline.ListPipelinesInput{}
	var pipelines []models.PipeLine
	err := client.ListPipelinesPages(params,
		func(resp *datapipeline.ListPipelinesOutput, lastPage bool) bool {
			for _, pipeline := range resp.PipelineIdList {
				pl := models.PipeLine{
					ID:   *pipeline.Id,
					Name: *pipeline.Name,
				}
				pipelines = append(pipelines, pl)
			}
			return true
		},
	)
	if err != nil {
		return pipelines, err
	}
	return pipelines, nil
}

func (p *DataPipeline) storePipelines(client *datapipeline.DataPipeline, pipelines []models.PipeLine) error {
	for _, pipeline := range pipelines {
		params := &datapipeline.GetPipelineDefinitionInput{
			PipelineId: &pipeline.ID,
		}
		resp, err := client.GetPipelineDefinition(params)
		if err != nil {
			return err
		}
		for _, po := range resp.PipelineObjects {
			for _, field := range po.Fields {
				if field.Key != nil && field.StringValue != nil {
					if *field.Key == "imageId" {
						pipeline.ImageId = *field.StringValue
						break
					}
					if *field.Key == "securityGroupIds" {
						resource := models.SecurityGroupResource{
							ResourceId:   pipeline.ID,
							Name:         pipeline.Name,
							ResourceType: "DataPipeline",
							AWSType:      "DataPipeline",
							GroupId:      *field.StringValue,
						}
						sgDb := p.DB.From("SecurityGroup")
						if err := sgDb.Save(&resource); err != nil {
							log.Fatal(err)
						}

					}
				}
			}
			if err := p.DB.Save(&pipeline); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func (p *DataPipeline) FetchAndStorePipelines(region string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error
	var pipelines []models.PipeLine
	session := session.New(&aws.Config{Region: &region})
	svc := datapipeline.New(session)
	if pipelines, err = p.listPipelines(svc); err != nil {
		return err
	}
	if storeErr := p.storePipelines(svc, pipelines); storeErr != nil {
		return err
	}
	return nil
}
