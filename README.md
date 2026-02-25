# steamcmd-healthcheck

A small Go program that can be used to send a [A2S_INFO](https://developer.valvesoftware.com/wiki/Server_queries#A2S_INFO) packet to a steamcmd gameserver. This can be used to check the health of your gameserver or as a healthcheck within a gameserver running in docker.

# Usage

This program will return 0 if it receives a successful response from your gameserver and otherwise return 1. Use the --verbose option to log 
output to your console.

```bash
./steamcmd-healthcheck heartbeat --host <gameserver ip> --port <gameserver port> --verbose
```

or with docker

```bash
docker run ghcr.io/radical-egg/steamcmd-healthcheck heartbeat --host <gameserver ip> --port <gameserver port> --verbose
```
# Building

```bash
go build -o steamcmd-healthcheck ./main.go
```

or with docker

```bash
docker build -t steamcmd-healthcheck .
```
