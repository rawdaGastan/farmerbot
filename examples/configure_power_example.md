# How to use configure power command

-   Get your redis DB address used in farmerbot
-   Create a new json file `config.json` and add your node options configurations:

```json
{
    "wakeUpThreshold": "<the threshold for resources usage that will need another node to be on, default is 80, optional>",
    "periodicWakeUp": "<daily time to wake up nodes for your farm, default is the time your run the command, format is 00:00AM or 00:00PM, optional>",
}
```

-   Then use the following code:

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

power, err := parser.ParseJSONIntoPower(jsonContent)
if err != nil {
    fmt.Print(err)
}

err = client.Call(ctx, "farmerbot.powermanager.Configure", []interface{}{power}, &err)
if err != nil {
    fmt.Print(err)
}
```
