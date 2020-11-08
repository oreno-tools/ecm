package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"strconv"
	"strings"
)

func GetClusters() []*string {
	res, err := svc.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}
	var clsNames []*string
	for _, r := range res.ClusterArns {
		splitedArn := strings.Split(*r, "/")
		clusterName := splitedArn[len(splitedArn)-1]
		clsNames = append(clsNames, aws.String(clusterName))
	}

	return clsNames
}

func GetClustersDetails(clsNames []*string, lcType string) ([][]string, []string) {
	input := &ecs.DescribeClustersInput{
		Clusters: clsNames,
	}

	res, err := svc.DescribeClusters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return nil, nil
	}
	var clsDatas [][]string
	for _, r := range res.Clusters {
		launchType := "EC2"
		if *r.RegisteredContainerInstancesCount == 0 && *r.RunningTasksCount > 0 {
			launchType = "FARGATE"
		} else if *r.RegisteredContainerInstancesCount == 0 && *r.RunningTasksCount == 0 {
			launchType = "---"
		}

		agentVersionStatus := "---"
		if launchType == "EC2" {
			agentVersionStatus = getClusterInstancesAgentVersionsStatus(*r.ClusterName)
		}

		clsData := []string{
			*r.ClusterName,
			launchType,
			strconv.FormatInt(*r.RegisteredContainerInstancesCount, 10),
			strconv.FormatInt(*r.RunningTasksCount, 10),
			strconv.FormatInt(*r.PendingTasksCount, 10),
			agentVersionStatus,
			*r.Status,
		}
		if lcType != "" {
			if lcType == launchType {
				clsDatas = append(clsDatas, clsData)
			}
		} else {
			clsDatas = append(clsDatas, clsData)
		}
	}
	clsHeader := []string{"Cluster Name", "Launch Type", "Container Instances", "Running Tasks", "Pending Tasks", "Agent Versions", "Status"}
	return clsDatas, clsHeader
}

func getClusterInstancesAgentVersionsStatus(clsName string) string {
	ids := getClusterInstances(clsName)
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clsName),
		ContainerInstances: ids,
	}

	result, err := svc.DescribeContainerInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}

	var agentVersions []*string
	for _, c := range result.ContainerInstances {
		// fmt.Println(c)
		// fmt.Println(*c.VersionInfo.AgentVersion)
		agentVersions = append(agentVersions, c.VersionInfo.AgentVersion)
	}

	for i := 1; i < len(agentVersions); i++ {
		if *agentVersions[i] != *agentVersions[0] {
			return "Mixed Version"
		} else {
			continue
		}
	}
	return *agentVersions[0]
}

func getClusterInstances(clsName string) []*string {
	input := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(clsName),
	}

	result, err := svc.ListContainerInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}

	var ids []*string
	for _, ca := range result.ContainerInstanceArns {
		cas := strings.Split(*ca, "/")
		id := cas[len(cas)-1]
		ids = append(ids, aws.String(id))
	}

	return ids
}

func GetClusterInstancesWithDatails(clsName string) ([][]string, []string) {
	ids := getClusterInstances(clsName)
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(*argCluster),
		ContainerInstances: ids,
	}

	result, err := svc.DescribeContainerInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return nil, nil
	}

	var insDatas [][]string
	for _, c := range result.ContainerInstances {
		// fmt.Println(c)
		cia := strings.Split(*c.ContainerInstanceArn, "/")
		id2 := cia[len(cia)-1]
		docker_version := strings.Replace(*c.VersionInfo.DockerVersion, "DockerVersion: ", "", 1)
		RunningTasksCount := strconv.FormatInt(*c.RunningTasksCount, 10)
		insData := []string{
			id2,
			*c.Ec2InstanceId,
			*c.VersionInfo.AgentVersion,
			docker_version,
			RunningTasksCount,
			*c.Status,
		}
		insDatas = append(insDatas, insData)
	}

	insHeader := []string{"Container Instnace", "Instance ID", "Agent Version", "Docker Version", "Running Tasks", "Status"}
	return insDatas, insHeader
}
