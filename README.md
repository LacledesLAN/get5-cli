# GET5-CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/LacledesLAN/get5-cli)](https://goreportcard.com/report/github.com/LacledesLAN/get5-cli)

`get5-cli` is an application for dynamically building `get5` configuration files from the command line.

[Get5](https://github.com/splewis/get5) is a standalone [SourceMod](http://www.sourcemod.net/) plugin for CS:GO servers for running matches. It is originally based on [pugsetup](https://github.com/splewis/csgo-pug-setup) and is inspired by [eBot](https://github.com/deStrO/eBot-CSGO). The core idea behind get5 is the ability to define all match details in a single config file; the main target use-case being tournaments and leagues (online or LAN). All that is required of the server-admins is to load match config file to the server and the match should run without any more manual actions from the admins.

## How It Works

`get5-cli` loads a "base" `get5` configuration file, combines it with values passed in via the CLI (command line interface), generates a new `get5` config file, and then launches a CSGO instance with `get5` using the generated config.

## Project Structure

* `/cmd/get5-cli` is the command line application.
* `/pkg/get5` is the base library for the program, intended to be imported by other programs such as [sourceseer](https://github.com/LacledesLAN/sourceseer).

## Motivation

At [Laclede's LAN](https://github.com/LacledesLAN/) we run the majority of our game servers in [Docker](https://www.docker.com/) fo reasons that are explained [here](https://github.com/LacledesLAN/README.1ST/blob/master/GameServers/DockerAndGameServers.md). To be able to containerize `get5` we need a way to dynamically inject get5 config files into the docker container.

## Why GO(lang) was Chosen

When choosing a language our key requirement was to able to compile native-binaries that could be added directly to Docker images without needing to include additional required dependencies. GO fit this criteria; and since we had multiple active projects being written in GO when this project was started it was a natural fit for our organization.

## Developer Notes

* This project includes configuration for [golangci-lint](https://github.com/golangci/golangci-lint); install the lint runner locally and then used the associated VSCode task to launch.
