# Codex Trunk
Trunk node for between chains

Trunk is part of the Codex muit-chain system. 
Some trunk node provides a decentralized state synchronization channel between the two chains. 
In particular, all chains in codex muit-chain system will have a channel by trunk nodes to codex.relay chain, so every chains can share their state to others by this.

- [Codex Trunk](#codex-trunk)
  - [Introduction](#introduction)
    - [Design](#design)
    - [Types](#types)
    - [Governance](#governance)
  - [Installation](#installation)
  - [Config](#config)


## Introduction

### Design

### Types

### Governance

## Installation

First build the trunk

```bash
go get -u -v github.com/codexnetwork/trunk
cd $GOPATH/src/github.com/codexnetwork/trunk
go get -v ./...
go install
```

You can use `-h` to get params for trunk:

```
trunk -h
Usage of ./trunk:
  -cfg string
    	config file path (default "./config.json")
  -d	run in debug mode
```

After make a config, can start trunk by:

```bash
trunk -cfg /path/to/config.json
```

## Config