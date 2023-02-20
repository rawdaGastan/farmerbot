# Farmerbot

Farmerbot is a service that a farmer can run allowing him to automatically manage the nodes of his farm.

// TODO: all commands examples

## How to start farmerbot

- Make sure to start redis server, and get db address for example <localhost:6379>

```bash
sudo systemctl start redis-server
```

- Create a new json file `config.json` and add your farm, nodes and power configurations:

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

- Get the binary

> Download the latest from the [releases page](https://github.com/rawdagastan/farmerbot/releases)

- Run the bot

After downloading the binary

```bash
sudo cp farmerbot /usr/local/bin
farmerbot -m <mnemonics> -c config.json -n dev -r <redis address> -d false
```

Where

- `<mnemonics>` is your farm mnemonics
- `config.json` is the json file of farmerbot configurations, with a default `config.json`
- `dev` is your network and can be main, qa and test with a default `dev`
- `<redis address>` is your redis DB address
- `false` is the value of debug mode with a default `false`

## Test

```bash
make test
```

## Release

- Check `goreleaser check`
- Create a tag `git tag -a v1.0.1 -m "release v1.0.1"`
- Push the tag `git push origin v1.0.1`
- A goreleaser workflow will release the tag
