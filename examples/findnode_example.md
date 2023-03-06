# How to use findnode command

- Get your redis DB address used in farmerbot
- Create a new json file `config.json` and add your node options configurations:

```json
{
    "certified": "<if you need a certified node, optional>",
    "dedicated": "<if you need a dedicated node, optional>",
    "publicConfig": "<if you need a publicConfig node, optional>",
    "publicIPs": "<number of public IPs you need, optional>",
    "capacity": {
        "SRU": "<enter needed sru, optional>",
        "MRU": "<enter needed mru, optional>",
        "HRU": "<enter needed hru, optional>",
        "CRU": "<enter needed cru, optional>"
    }
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

nodeOptions, err := parser.ParseJSONIntoNodeOptions(jsonContent)
if err != nil {
    fmt.Print(err)
}

var node uint32
err = client.Call(ctx, "farmerbot.nodemanager.FindNode", []interface{}{nodeOptions, []uint{}}, &node)
if err != nil {
    fmt.Print(err)
}
```
