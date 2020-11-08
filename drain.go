package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
)

func ExecuteDrainAll(clusterName string, agentVersion string) bool {
	insDatas, _ := GetClusterInstancesWithDatails(clusterName)

	var drainTargets [][]string
	for _, insd := range insDatas {
		if !strings.EqualFold(insd[2], agentVersion) {
			drainTargets = append(drainTargets, insd)
		}
	}

	if len(drainTargets) == 0 {
		fmt.Println("The drain target instance does not exist.")
		return false
	}

	if checkDrainCondition(insDatas, drainTargets) {
		result := executeDrain(clusterName, drainTargets)
		if result {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func ExecuteDrain(clusterName string, instanceId string) bool {
	insDatas, _ := GetClusterInstancesWithDatails(clusterName)

	var drainTargets [][]string
	drainTargets = append(drainTargets, []string{instanceId})

	if checkDrainCondition(insDatas, drainTargets) {
		result := executeDrain(clusterName, drainTargets)
		if result {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func checkDrainCondition(insDatas [][]string, drainTargets [][]string) bool {
	// 全体のインスタンス数が drain 対象のインスタンス数よりも多いことが条件
	if len(insDatas) >= (len(drainTargets) * 2) {
		fmt.Printf("Cluster Instances: %d\nDrain Target Instances: %d\n", len(insDatas), len(drainTargets))
		return true
	} else {
		fmt.Printf("There are not enough instances in the cluster.\n")
		return false
	}
}

func executeDrain(clusterName string, drainTargets [][]string) bool {
	fmt.Printf("Do you want to continue processing? (y/n): ")
	var stdin string
	fmt.Scan(&stdin)
	switch stdin {
	case "y", "Y":
		var conInsts []*string
		for _, conIns := range drainTargets {
			conInsts = append(conInsts, aws.String(conIns[0]))
		}

		input := &ecs.UpdateContainerInstancesStateInput{
			Cluster:            aws.String(clusterName),
			Status:             aws.String("DRAINING"),
			ContainerInstances: conInsts,
		}

		insState, err := svc.UpdateContainerInstancesState(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				fmt.Println(aerr.Error())
			} else {
				fmt.Println(err.Error())
			}
			return false
		}

		for _, r := range insState.ContainerInstances {
			st := fmt.Sprintf("Instance ID: %s Status: %s", *r.Ec2InstanceId, *r.Status)
			fmt.Println(st)
		}
		return true
	case "n", "N":
		fmt.Println("Interrupted.")
		return true
	default:
		fmt.Println("Interrupted.")
		return true
	}
}
