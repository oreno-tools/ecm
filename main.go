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
	AppVersion = "0.0.2"
)

var (
	argVersion      = flag.Bool("version", false, "Print version number.")
	argCluster      = flag.String("cluster", "", "Set a AutoScaling Group Name.")
	argDrain        = flag.Bool("drain", false, "Execute draining.")
	argDrainAll     = flag.Bool("drain-all", false, "Execute all instance draining.")
	argAgentVersion = flag.String("agent-version", "", "Specify the agent version.")
	argInstance     = flag.String("instance", "", "Specify the instances.")

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
			RunningTasksCount,
			*c.Status,
		}
		insDatas = append(insDatas, insData)
	}

	insHeader := []string{"Container Instnace", "Instance ID", "Agent Version", "Docker Version", "Running Tasks", "Status"}
	printTable(insDatas, insHeader)
	// os.Exit(0)

	if *argInstance == "" && !*argDrainAll {
		os.Exit(0)
	}

	// fmt.Println(insDatas)
	var drainTargets [][]string
	if *argDrainAll {
		for _, insd := range insDatas {
			if !strings.EqualFold(insd[2], *argAgentVersion) {
				drainTargets = append(drainTargets, insd)
			}
		}
	} else if *argInstance != "" {
		drainTargets = append(drainTargets, []string{*argInstance})
	}

	if len(drainTargets) == 0 {
		fmt.Println("The drain target instance does not exist.")
		os.Exit(1)
	}

	// 全体のインスタンス数が drain 対象のインスタンス数よりも多いことが条件
	if len(insDatas) >= (len(drainTargets) * 2) {
		fmt.Printf("Cluster Instances: %d\nDrain Target Instances: %d\n", len(insDatas), len(drainTargets))
	} else {
		fmt.Printf("There are not enough instances in the cluster.\n")
		os.Exit(1)
	}

	fmt.Printf("Do you want to continue processing? (y/n): ")
	var stdin string
	fmt.Scan(&stdin)
	switch stdin {
	case "y", "Y":
		var conInsts []*string
		for _, conIns := range drainTargets {
			conInsts = append(conInsts, aws.String(conIns[0]))
		}

		////////////////////////////////////////////////////////////////////////
		input3 := &ecs.UpdateContainerInstancesStateInput{
			Cluster:            aws.String(*argCluster),
			Status:             aws.String("DRAINING"),
			ContainerInstances: conInsts,
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
	case "n", "N":
		fmt.Println("Interrupted.")
		os.Exit(0)
	default:
		fmt.Println("Interrupted.")
		os.Exit(0)
	}

}
