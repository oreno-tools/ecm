package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	_ "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	_ "github.com/gosuri/uilive"
	_ "github.com/kyokomi/emoji"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
	_ "time"
)

const (
	AppVersion = "0.0.1"
)

var (
	argVersion  = flag.Bool("version", false, "Print version number.")
	argCluster  = flag.String("cluster", "", "Set a AutoScaling Group Name.")
	argDrain    = flag.Bool("drain", false, "Execute draining.")
	argInstance = flag.String("instance", "", "Specify the instances.")
	// argStatus   = flag.String("status", "", "Specify the instances.")

	svc = ecs.New(session.New())
)

func printTable(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}

func main() {
	flag.Parse()

	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	if *argCluster == "" {
		input := &ecs.ListClustersInput{}
		result, err := svc.ListClusters(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				fmt.Println(aerr.Error())
			} else {
				fmt.Println(err.Error())
			}
			return
		}
		var clsDatas [][]string
		for _, r := range result.ClusterArns {
			// fmt.Println(*r)
			splitedArn := strings.Split(*r, "/")
			clusterName := splitedArn[len(splitedArn)-1]
			clsData := []string{
				clusterName,
			}
			clsDatas = append(clsDatas, clsData)
		}
		clsHeader := []string{"Cluster Name"}
		printTable(clsDatas, clsHeader)
		os.Exit(0)
	}

	////////////////////////////////////////////////////////////////////////
	input := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(*argCluster),
	}

	result, err := svc.ListContainerInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	var ids []*string
	for _, ca := range result.ContainerInstanceArns {
		cas := strings.Split(*ca, "/")
		id := cas[len(cas)-1]
		ids = append(ids, aws.String(id))
	}

	////////////////////////////////////////////////////////////////////////
	input2 := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(*argCluster),
		ContainerInstances: ids,
	}

	result2, err := svc.DescribeContainerInstances(input2)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	// fmt.Println(result2)
	var insDatas [][]string
	for _, c := range result2.ContainerInstances {
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
			RunningTasksCount, *c.Status,
		}
		insDatas = append(insDatas, insData)
	}

	insHeader := []string{"Container Instnace", "Instance ID", "Agent Version", "Docker Version", "Running Tasks", "Status"}
	printTable(insDatas, insHeader)
	// os.Exit(0)

	if *argInstance == "" {
		os.Exit(0)
	}

	////////////////////////////////////////////////////////////////////////
	input3 := &ecs.UpdateContainerInstancesStateInput{
		Cluster: aws.String(*argCluster),
		Status:  aws.String("DRAINING"),
		ContainerInstances: []*string{
			aws.String(*argInstance),
		},
	}

	result3, err := svc.UpdateContainerInstancesState(input3)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	for _, r3 := range result3.ContainerInstances {
		fmt.Println(*r3.Status)
	}
}
