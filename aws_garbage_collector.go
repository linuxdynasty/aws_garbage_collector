package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/linuxdynasty/aws_garbage_collector/controllers"
	"github.com/linuxdynasty/aws_garbage_collector/ec2"
	"github.com/linuxdynasty/aws_garbage_collector/iam"
	"github.com/linuxdynasty/aws_garbage_collector/shared"
)

func init() {
	shared.DBC = shared.PrepareDb("my.db")
}

func processInUse() {
	ec2DB := ec2.DB(shared.DBC)
	log.Printf("Processing In Use AMIs ")
	ec2DB.EC2ProcessInUse()
	log.Printf("Done Processing In Use AMIs ")

	log.Printf("Processing In Use SecurityGroups ")
	ec2DB.SGProcessInUse()
	log.Printf("Done Processing In Use SecurityGroups ")

	iamDB := iam.DB(shared.DBC)
	log.Printf("Processing In Use IAM ")
	iamDB.InstanceProfileProcessInUse()
	log.Printf("Done Processing In Use IAM ")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	shared.DefaultRegion = "us-west-2"
	regions := shared.Regions(shared.DefaultRegion)
	var wg sync.WaitGroup
	wg.Add(1)
	go ProcessIAM(shared.DBC, shared.DefaultRegion, &wg)
	wg.Add(len(regions) * 6)
	for _, region := range regions {
		go ProcessDataPipelines(shared.DBC, region, &wg)
		go ProcessLaunchConfigurations(shared.DBC, region, &wg)
		go ProcessEC2(shared.DBC, region, &wg)
		go ProcessElastiCache(shared.DBC, region, &wg)
		go ProcessRDS(shared.DBC, region, &wg)
		go ProcessRedShift(shared.DBC, region, &wg)
	}
	wg.Wait()
	processInUse()
	log.Println("Starting Web Server")
	mux := httprouter.New()
	mux.GET("/amis", controllers.GetAMIs)
	mux.GET("/amis/:region", controllers.GetAMIs)
	mux.POST("/amis/:region", controllers.DeleteAMIs)
	mux.GET("/iam/policies", controllers.GetPolicies)
	mux.GET("/iam/policies/:region", controllers.GetPolicies)
	mux.GET("/iam/profiles", controllers.GetInstanceProfiles)
	mux.GET("/iam/profiles/:region", controllers.GetInstanceProfiles)
	mux.GET("/securitygroups", controllers.GetSecurityGroups)
	mux.GET("/securitygroups/:region", controllers.GetSecurityGroups)
	mux.POST("/securitygroups/:region", controllers.DeleteSecurityGroups)
	mux.GET("/launchconfigurations", controllers.GetLaunchConfigurations)
	mux.GET("/launchconfigurations/:region", controllers.GetLaunchConfigurations)
	mux.POST("/launchconfigurations/:region", controllers.DeleteLaunchConfigurations)
	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}
	server.ListenAndServe()
}
