package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/linuxdynasty/aws_garbage_collector/api"
	"github.com/linuxdynasty/aws_garbage_collector/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func GetAMIs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var inUse string
	var inUseByDataPipeLine string
	var inUseByInstance string
	var inUseByLC string
	var all bool = false

	region = p.ByName("region")
	if region != "" {
		response, regionIsValid := parseRegion(region)
		if !regionIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	r.ParseForm()
	inUse = r.FormValue("in_use")
	if inUse != "" {
		response, valueIsValid := parseFormStringBool(inUse)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByInstance = r.FormValue("in_use_by_instance")
	if inUseByInstance != "" {
		response, valueIsValid := parseFormStringBool(inUseByInstance)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByDataPipeLine = r.FormValue("in_use_by_datapipeline")
	if inUseByDataPipeLine != "" {
		response, valueIsValid := parseFormStringBool(inUseByDataPipeLine)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByLC = r.FormValue("in_use_by_lc")
	if inUseByLC != "" {
		response, valueIsValid := parseFormStringBool(inUseByLC)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}

	if inUse == "" && inUseByLC == "" && inUseByDataPipeLine == "" && inUseByInstance == "" {
		all = true
	}
	ec2Data := ec2.DB(shared.DBC)

	response := api.GetAllAMIs(&ec2Data, region, inUse, inUseByLC, inUseByInstance, inUseByDataPipeLine, all)
	marshalled, _ = json.Marshal(response)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))

}

func DeleteAMIs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var unUsedIds []string
	var response models.DeleteResponseApi
	var statuses []models.DeleteStatus

	ec2DB := ec2.DB(shared.DBC)

	if p.ByName("region") != "" {
		region = p.ByName("region")
		if !shared.IsRegionValid(region, shared.DefaultRegion) {
			response.Count = 0
			response.Message = fmt.Sprintf("Invalid value: %s, Must be a valid AWS Region", region)
			marshalled, _ = json.Marshal(response)
			pretty, _ := shared.PrettyJSON(marshalled)
			fmt.Fprintf(w, string(pretty))
			return
		}
	} else {
		response.Count = 0
		response.Message = fmt.Sprintf("Invalid URL: %s: Must be /amis/region_name_goes_here", r.URL.RequestURI)
		marshalled, _ = json.Marshal(response)
		pretty, _ := shared.PrettyJSON(marshalled)
		fmt.Fprintf(w, string(pretty))
		return
	}
	r.ParseForm()
	amiIds := r.Form["id"]
	allUnused := r.FormValue("all_unused")
	if allUnused == "true" {
		unUsedIds = append(unUsedIds, api.GetAllUnusedAMIIds(&ec2DB, region)...)
	}
	if len(amiIds) > 0 && allUnused == "" {
		statuses = ec2DB.DeleteByIds(amiIds, region)
	} else if allUnused != "" {
		statuses = ec2DB.DeleteByIds(unUsedIds, region)
	} else {
		response.Count = 0
		response.Message = fmt.Sprint("Must be at least one AMI ID")
	}
	response.Count = len(statuses)
	response.Data = statuses
	marshalled, _ = json.Marshal(response)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))
	return
}
