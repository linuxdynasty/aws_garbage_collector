package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/linuxdynasty/aws_garbage_collector/api"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func GetPolicies(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var inUse string
	var inUseByRoles string
	var inUseByUsers string
	var inUseByGroups string
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
	inUseByRoles = r.FormValue("in_use_by_roles")
	if inUseByRoles != "" {
		response, valueIsValid := parseFormStringBool(inUseByRoles)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByGroups = r.FormValue("in_use_by_groups")
	if inUseByGroups != "" {
		response, valueIsValid := parseFormStringBool(inUseByGroups)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByUsers = r.FormValue("in_use_by_users")
	if inUseByUsers != "" {
		response, valueIsValid := parseFormStringBool(inUseByUsers)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	if inUse == "" && inUseByRoles == "" && inUseByUsers == "" && inUseByGroups == "" {
		all = true
	}
	iamData := iam.DB(shared.DBC)

	response := api.GetAllPolicies(&iamData, region, inUse, inUseByRoles, inUseByUsers, inUseByGroups, all)
	marshalled, _ = json.Marshal(response)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))
}

func GetInstanceProfiles(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var inUse string
	var inUseByRoles string
	var inUseByInstances string
	var inUseByLCs string
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
	inUseByRoles = r.FormValue("in_use_by_roles")
	if inUseByRoles != "" {
		response, valueIsValid := parseFormStringBool(inUseByRoles)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByInstances = r.FormValue("in_use_by_instances")
	if inUseByInstances != "" {
		response, valueIsValid := parseFormStringBool(inUseByInstances)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	inUseByLCs = r.FormValue("in_use_by_launch_configurations")
	if inUseByLCs != "" {
		response, valueIsValid := parseFormStringBool(inUseByLCs)
		if !valueIsValid {
			fmt.Fprintf(w, string(response))
			return
		}
	}
	if inUse == "" && inUseByRoles == "" && inUseByInstances == "" && inUseByLCs == "" {
		all = true
	}
	iamData := iam.DB(shared.DBC)

	response := api.GetAllInstanceProfiles(&iamData, region, inUse, inUseByRoles, inUseByInstances, inUseByLCs, all)
	marshalled, _ = json.Marshal(response)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))
}
