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

- Run:

```bash
farmerbot farmmanager define -c config.json -m <mnemonics> -n dev -r <redis address> -d false -l farmerbot.log
```

Where:

- `-c config.json` is the json file of farmerbot configurations, with a default `config.json`.
- `-m <mnemonics>` is your farm mnemonics.
- `-n dev` is your network and can be main, qa and test with a default `dev`.
- `-r <redis address>` is your redis DB address.
- `-d false` is the value of debug mode with a default `false`.
- `-l farmerbot.log` is log file to include logs generated by farmerbot with a default `farmerbot.log`.