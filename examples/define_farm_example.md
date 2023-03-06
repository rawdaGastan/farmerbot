# How to use define farm command

- Get your redis DB address used in farmerbot
- Create a new json file `config.json` and add your node options configurations:

```json
{
    "id": "<your farm ID, required>",
    "description": "<farm description, optional>",
    "publicIPs": "<number of public ips in your farm, optional>"
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

farm, err := parser.ParseJSONIntoFarm(jsonContent)
if err != nil {
    fmt.Print(err)
}

err = client.Call(ctx, "farmerbot.farmmanager.Define", []interface{}{farm}, &err)
if err != nil {
    fmt.Print(err)
}
```
