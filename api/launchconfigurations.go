package api

import (
	"github.com/asdine/storm/q"
	"github.com/linuxdynasty/aws_garbage_collector/asgs"
	"github.com/linuxdynasty/aws_garbage_collector/models"
)

func getAllAsgs(l *asgs.LC, lcs []models.LaunchConfiguration) []models.LaunchConfigDetails {
	var lcsDetails []models.LaunchConfigDetails

	for _, lc := range lcs {
		lcsDetail := models.LaunchConfigDetails{
			ID:              lc.ID,
			Name:            lc.Name,
			InUse:           lc.InUse,
			ImageId:         lc.ImageId,
			InstanceProfile: lc.InstanceProfile,
			Region:          lc.Region,
		}
		var asgs []models.AutoScaleGroup

		l.ASGBucket.Find("LaunchConfigurationName", lc.Name, &asgs)
		lcsDetail.AutoScaleGroups = asgs

		lcsDetails = append(lcsDetails, lcsDetail)
	}
	return lcsDetails
}

func GetAllLaunchConfigurations(l *asgs.LC, region, inUse string, all bool) *models.GetResponseApi {
	var lcs []models.LaunchConfiguration
	if region != "" {
		if all {
			l.DB.Find("Region", region, &lcs)
		} else {
			l.DB.Select(q.And(
				q.Eq("InUse", inUse),
				q.Eq("Region", region),
			)).Find(&lcs)
		}
	} else {
		if all {
			l.DB.All(&lcs)
		} else {
			l.DB.Find("InUse", inUse, &lcs)
		}
	}
	lcsDetails := getAllAsgs(l, lcs)

	apiResponse := models.GetResponseApi{
		Data:  lcsDetails,
		Count: len(lcsDetails),
	}
	return &apiResponse
}
