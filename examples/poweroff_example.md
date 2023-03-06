# How to use poweroff command

- Get your redis DB address used in farmerbot
- Get your desired node for example: 1
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

nodeID := uint32(1)
err = client.Call(ctx, "farmerbot.powermanager.PowerOff", []interface{}{nodeID}, &err)
if err != nil {
    fmt.Print(err)
}
```
