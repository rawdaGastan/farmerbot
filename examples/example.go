// Package main
package main

import (
	"context"
	"fmt"

	"github.com/rawdaGastan/farmerbot/client"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/threefoldtech/zbus"
)

func examples(ctx context.Context, redisAddr string) error {
	address := fmt.Sprintf("tcp://%s", redisAddr)
	zBusClient, err := zbus.NewRedisClient(address)
	if err != nil {
		return err
	}

	client := client.NewFarmerClient(zBusClient)

	err = client.Call(ctx, "farmerbot.farmmanager.Define", []interface{}{models.Farm{ID: 1}}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	err = client.Call(ctx, "farmerbot.nodemanager.Define", []interface{}{models.Node{ID: 1}}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	var node uint32
	err = client.Call(ctx, "farmerbot.nodemanager.FindNode", []interface{}{models.NodeOptions{}, []uint{}}, &node)
	fmt.Printf("node ID: %v\n", node)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	err = client.Call(ctx, "farmerbot.powermanager.Configure", []interface{}{models.Power{}}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	err = client.Call(ctx, "farmerbot.powermanager.PowerOn", []interface{}{uint32(1)}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	err = client.Call(ctx, "farmerbot.powermanager.PowerOff", []interface{}{uint32(1)}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	err = client.Call(ctx, "farmerbot.powermanager.PeriodicWakeup", []interface{}{}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	err = client.Call(ctx, "farmerbot.powermanager.PowerManagement", []interface{}{}, &err)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	return nil
}

func main() {
	redisAddr := "localhost:6379"
	err := examples(context.Background(), redisAddr)
	if err != nil {
		fmt.Println("examples failed with error: ", err)
	}
}
