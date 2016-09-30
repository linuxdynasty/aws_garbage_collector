package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func parseRegion(region string) ([]byte, bool) {
	var returnOutput []byte
	var regionIsValid bool

	if !shared.IsRegionValid(region, shared.DefaultRegion) {
		response := models.GetResponseApi{
			Count:   0,
			Message: fmt.Sprintf("Invalid value: %s, Must be a valid AWS Region", region),
		}
		marshalled, _ := json.Marshal(response)
		returnOutput, _ = shared.PrettyJSON(marshalled)
		regionIsValid = false
	} else {
		regionIsValid = true
	}

	return returnOutput, regionIsValid
}
func parseFormStringBool(value string) ([]byte, bool) {
	var returnOutput []byte
	var valueIsValid bool

	if value != "true" && value != "false" {
		response := models.GetResponseApi{
			Count:   0,
			Message: fmt.Sprintf("Invalid value: %s, Must pass either true or false", value),
		}
		marshalled, _ := json.Marshal(response)
		returnOutput, _ = shared.PrettyJSON(marshalled)
	} else {
		valueIsValid = true
	}

	return returnOutput, valueIsValid
}
