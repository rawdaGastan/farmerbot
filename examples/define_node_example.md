# How to use define node command

- Get your redis DB address used in farmerbot
- Create a new json file `config.json` and add your node configurations (some are required):

```json
{
    "id": "<your node ID, required>",
    "twinID": "<your node twin ID, required>",
    "farmID": "<your node farm ID, optional>",
    "farmID": "<your node farm ID, optional>",
    "description": "<description, optional>",
    "certified": "<if node is certified, optional>",
    "dedicated": "<if node is dedicated, optional>",
    "publicConfig": "<if node has public config, optional>",
    "publicIPsUsed": "<number of node used public ips, optional>",
    "hasActiveRentContract": "<if node has an active rent contract, optional>",
    "wgPorts": "<list of node wireguard ports, optional>",
    "resources": {
        "overProvisionCPU": "<how much node allow over provisioning the CPU , default is 1, range: [1;4], optional>",
        "total": {
            "SRU": "<node SRU>, required",
            "MRU": "<node MRU>, required",
            "HRU": "<node HRU>, required",
            "CRU": "<node CRU>, required",
        }
    },
    "powerState": {
        "on": "<if node power state is on, default is true, optional>",
        "wakingUp": "<if node power state is waking up, optional>",
        "off": "<if node power state is off, optional>",
        "shuttingDown": "<if node power state is shutting down, optional>"
    },
    "timeoutClaimedResources": "<timeout to update claiming resources from node, default is after 30 minutes, optional>",
    "lastTimePowerStateChanged": "<last time node power changed, optional>",
    "lastTimeAwake": "<last time node was waking up, optional>",
}
```

- Then use the following code:

```go
// Package main
package main

import (
    "context"
    "fmt"   

    "github.com/rawdaGastan/farmerbot/client"
    "github.com/rawdaGastan/farmerbot/internal/models"
    "github.com/threefoldtech/zbus"
)

address := fmt.Sprintf("tcp://%s", redisAddr)
zBusClient, err := zbus.NewRedisClient(address)
if err != nil {
    return err
}

client := client.NewFarmerClient(zBusClient)

jsonContent, err := parser.ReadFile("config.json")
if err != nil {
    fmt.Print(err)
}

node, err := parser.ParseJSONIntoNode(jsonContent)
if err != nil {
    fmt.Print(err)
}

err = client.Call(ctx, "farmerbot.nodemanager.Define", []interface{}{node}, &err)
if err != nil {
    fmt.Print(err)
}
```
