# Farmerbot

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/dc1cd40b31324ff1b80d9706b08837e8)](https://www.codacy.com/gh/rawdaGastan/farmerbot/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=rawdaGastan/farmerbot&amp;utm_campaign=Badge_Grade) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-79%25-brightgreen.svg?longCache=true&style=flat)</a>

Farmerbot is a service that a farmer can run allowing him to automatically manage the nodes of his farm.

## How to start farmerbot

-   Make sure to start redis server, and get redis DB address for example: <localhost:6379>

```bash
sudo systemctl start redis-server
```

-   Create a new json file `config.json` and add your farm, nodes and power configurations:

```json
{
    "farm": {
        "id": "<your farm ID>"
    },
    "nodes": [{
        "id": "<your node ID>",
        "twinID": "<your node twin ID>",
        "resources": {
            "total": {
                "SRU": "<enter total sru>",
                "MRU": "<enter total mru>",
                "HRU": "<enter total hru>",
                "CRU": "<enter total cru>"
            }
        }
    }],
    "power": {
        "periodicWakeUp": "08:30AM",
        "wakeUpThreshold": 80
    }
}
```

-   Get the binary

> Download the latest from the [releases page](https://github.com/rawdagastan/farmerbot/releases)

-   Run the bot

After downloading the binary

```bash
sudo cp farmerbot /usr/local/bin
farmerbot -c config.json -m <mnemonics> -n dev -r <redis address> -d false -l farmerbot.log
```

Where:

-   `-c config.json` is the json file of farmerbot configurations, with a default `config.json`.
-   `-m <mnemonics>` is your farm mnemonics.
-   `-n dev` is your network and can be main, qa and test with a default `dev`.
-   `-r <redis address>` is your redis DB address.
-   `-d false` is the value of debug mode with a default `false`.
-   `-l farmerbot.log` is log file to include logs generated by farmerbot with a default `farmerbot.log`.

> Note: **`30 minutes`** are set for a timeout node power change

## Server

You can start farmerbot server with the following command

```bash
farmerbot server -m <mnemonics> -n <grid network> -r <redis address> -d <debug> -l <log file>
```

## Supported commands

-   farmerbot powermanager [configure](/examples/configure_power_example.md)
-   farmerbot nodemanager [define](/examples/define_node_example.md)
-   farmerbot farmmanager [define](/examples/define_farm_example.md)

-   farmerbot powermanager [poweron](/examples/poweron_example.md)
-   farmerbot powermanager [poweroff](/examples/poweroff_example.md)
-   farmerbot nodemanager [findnode](/examples/findnode_example.md)

For more examples and explanations for supported commands, see the [examples](/examples)

## Examples

To run examples:

-   Run the [server](#server) then:

```bash
go run examples/example.go
```

## Version

You can get the latest version of the farmerbot by running the following command:

```bash
farmerbot version
```

## Test

```bash
make test
```

## Release

-   Check `goreleaser check`
-   Create a tag `git tag -a v1.0.1 -m "release v1.0.1"`
-   Push the tag `git push origin v1.0.1`
-   A goreleaser workflow will release the created tag.
