package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

const (
	AppVersion = "0.0.4"
)

var (
	argVersion      = flag.Bool("version", false, "Print version number.")
	argCluster      = flag.String("cluster", "", "Set a AutoScaling Group Name.")
	argDrain        = flag.Bool("drain", false, "Execute draining.")
	argDrainAll     = flag.Bool("drain-all", false, "Execute all instance draining.")
	argAgentVersion = flag.String("agent-version", "", "Specify the agent version.")
	argInstance     = flag.String("instance", "", "Specify the instances.")
	argLaunchType   = flag.String("type", "", "Specify the launch type (EC2 or Fargate)")
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
		clsNames := GetClusters()
		clsDatas, clsHeader := GetClustersDetails(clsNames, *argLaunchType)
		if clsDatas != nil {
			printTable(clsDatas, clsHeader)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		if *argDrain && *argInstance != "" {
			result := ExecuteDrain(*argCluster, *argInstance)
			if result {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
		//
		if *argDrainAll && *argAgentVersion != "" {
			result := ExecuteDrainAll(*argCluster, *argAgentVersion)
			if result {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
		//
		insDatas, insHeader := GetClusterInstancesWithDatails(*argCluster)
		if insDatas != nil {
			printTable(insDatas, insHeader)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
