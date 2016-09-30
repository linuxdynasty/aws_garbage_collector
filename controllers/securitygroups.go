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

func GetSecurityGroups(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	var inUse string
	var inUseBySgOnly string
	var all bool = false

	if p.ByName("region") != "" {
		region = p.ByName("region")
		if !shared.IsRegionValid(region, shared.DefaultRegion) {
			sgs := models.GetResponseApi{
				Count:   0,
				Message: fmt.Sprintf("Invalid value: %s, Must be a valid AWS Region", region),
			}
			marshalled, _ = json.Marshal(sgs)
			pretty, _ := shared.PrettyJSON(marshalled)
			fmt.Fprintf(w, string(pretty))
			return
		}
	}

	r.ParseForm()
	if val := r.FormValue("in_use"); val != "" {
		if val != "true" && val != "false" {
			sgs := models.GetResponseApi{
				Count:   0,
				Message: fmt.Sprintf("Invalid value: %s, Must pass either true or false", val),
			}
			marshalled, _ = json.Marshal(sgs)
			pretty, _ := shared.PrettyJSON(marshalled)
			fmt.Fprintf(w, string(pretty))
			return
		}
		inUse = val
	}
	r.ParseForm()
	if val := r.FormValue("in_use_by_sg_only"); val != "" {
		if val != "true" && val != "false" {
			sgs := models.GetResponseApi{
				Count:   0,
				Message: fmt.Sprintf("Invalid value: %s, Must pass either true or false", val),
			}
			marshalled, _ = json.Marshal(sgs)
			pretty, _ := shared.PrettyJSON(marshalled)
			fmt.Fprintf(w, string(pretty))
			return
		}
		inUseBySgOnly = val
	}

	if r.FormValue("in_use") == "" && r.FormValue("in_use_by_sg_only") == "" {
		all = true
	}
	group := ec2.DB(shared.DBC)

	sgs := api.GetAllSecurityGroups(&group, region, inUse, inUseBySgOnly, all)
	marshalled, _ = json.Marshal(sgs)
	pretty, _ := shared.PrettyJSON(marshalled)
	fmt.Fprintf(w, string(pretty))
}

func DeleteSecurityGroups(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var region string
	var marshalled []byte
	group := ec2.DB(shared.DBC)

	if p.ByName("region") != "" {
		region = p.ByName("region")
		if !shared.IsRegionValid(region, shared.DefaultRegion) {
			sgs := models.DeleteResponseApi{
				Count:   0,
				Message: fmt.Sprintf("Invalid value: %s, Must be a valid AWS Region", region),
			}
			marshalled, _ = json.Marshal(sgs)
			pretty, _ := shared.PrettyJSON(marshalled)
			fmt.Fprintf(w, string(pretty))
			return
		}
	}
	r.ParseForm()
	ids := r.Form["id"]
	if len(ids) > 0 {
		statuses := group.SGDeleteByIds(ids, region)
		marshalled, _ = json.Marshal(statuses)
		pretty, _ := shared.PrettyJSON(marshalled)
		fmt.Fprintf(w, string(pretty))
	} else {
		sgs := models.DeleteResponseApi{
			Count:   0,
			Message: fmt.Sprint("Must be at least one security group id"),
		}
		marshalled, _ = json.Marshal(sgs)
		pretty, _ := shared.PrettyJSON(marshalled)
		fmt.Fprintf(w, string(pretty))
		return
	}
}
