# SRB2Kart Ganymede - Discord bot

A Discord bot for displaying SRB2Kart server status, written in `go` this time :p

Can be (for the most part) a drop-in alternative for [dione](https://github.com/Indev450/dione)

## Installing

```sh
go install github.com/Indev450/ganymede@latest
```

This should install latest version of bot into your go path (usually `~/go/bin/`)

## Running

Assuming program is installed into `~/go/bin/`

```sh
DISCORD_TOKEN=<your token> SRB2KART_ADDRESS=<ip:port of server> ~/go/bin/ganymede
```

## Environment variables

`DISCORD_TOKEN` - required, bot token

`DISCORD_CMDPREFIX` - optional, will be added to all commands (for example, a bot for ring racers server might use `DISCORD_CMDPREFIX="rr-"`, resulting into
commands being `/rr-players` and `/rr-gamemode`)

`SRB2KART_ADDRESS` - optional, address of srb2kart server, which would be asked for info (format is `host:port`. host defaults to `127.0.0.1`, port defaults to `5029`)

`SRB2KART_PROTO` - optional, and defaults to `srb2kart-16p` (which is vanilla kart server). Other currently supported options are:

- `ringracers-16p` (not sure if there are non-16p clients, so it is with 16p suffix for now)
- `blankart` (blankart is in active dev, so support may or may not break occasionally :p)

`SRB2KART_GAMEMODEFILE` - optional, path to file which would be read to fetch server gamemodes. File should store each gamemode on new line. If this
variable is not found, /gamemode command will not be available

`SRB2KART_SEEDPLAYER` - optional, don't count player with that name in status
