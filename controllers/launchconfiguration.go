package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/linuxdynasty/aws_garbage_collector/api"
	"github.com/linuxdynasty/aws_garbage_collector/asgs"
	"github.com/linuxdynasty/aws_garbage_collector/models"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func GetLaunchConfigurations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var inUse string
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

	if inUse == "" {
		all = true
	}
	lc := asgs.DB(shared.DBC)

	lcs := api.GetAllLaunchConfigurations(&lc, region, inUse, all)
	marshalled, _ = json.Marshal(lcs)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))
}

func DeleteLaunchConfigurations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var unUsedNames []string
	var response models.DeleteResponseApi
	var statuses []models.DeleteStatus

	launchConfig := asgs.DB(shared.DBC)

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
	}
	r.ParseForm()
	names := r.Form["name"]
	allUnused := r.FormValue("all_unused")
	if allUnused == "true" {
		var lcs []models.LaunchConfiguration
		if searchErr := launchConfig.DB.Find("InUse", "false", &lcs); searchErr == nil {
			for _, lc := range lcs {
				unUsedNames = append(unUsedNames, lc.Name)
			}
		}
	}
	if len(names) > 0 && allUnused == "" {
		statuses = launchConfig.DeleteByName(names, region)
	} else if allUnused != "" {
		statuses = launchConfig.DeleteByName(unUsedNames, region)
	} else {
		response.Count = 0
		response.Message = fmt.Sprint("Must be at least one launch configuration name")
	}
	response.Count = len(statuses)
	response.Data = statuses
	marshalled, _ = json.Marshal(response)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))
	return
}
