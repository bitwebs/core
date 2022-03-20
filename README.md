<p>&nbsp;</p>
<p align="center">

<img src="core_logo.svg" width=500>

</p>

<p align="center">
Full-node software implementing the IQ protocol<br/><br/>

<a href="https://codecov.io/gh/bitwebs/core">
    <img src="https://codecov.io/gh/bitwebs/core/branch/develop/graph/badge.svg">
</a>
<a href="https://goreportcard.com/report/github.com/bitwebs/core">
    <img src="https://goreportcard.com/badge/github.com/bitwebs/core">
</a>

</p>

<p align="center">
  <a href="https://docs.iqchain.network/"><strong>Explore the Docs »</strong></a>
  <br />
  <br/>
  <a href="https://docs.iqchain.network/docs/develop/module-specifications/README.html">IQ Core reference</a>
  ·
  <a href="https://pkg.go.dev/github.com/bitwebs/core?tab=subdirectories">Go API</a>
  ·
  <a href="https://api.iqchain.network/swagger/#/">Rest API</a>
  ·
  <a href="https://github.com/bitwebs/iq.py">Python SDK</a>
  ·
  <a href="https://bitwebs.github.io/iq.js/">IQChain.js</a>
  ·
  <a href="https://finder.iqchain.network/">Finder</a>
  ·
  <a href="https://station.iqchain.network/">Station</a>
</p>

<br/>

## Table of Contents <!-- omit in toc -->

- [What is IQ Chain?](#what-is-iq-chain)
- [Installation](#installation)
  - [Binaries](#binaries)
  - [From Source](#from-source)
- [`IQd`](#iqd)
- [Node Setup](#node-setup)
  - [Join the mainnet](#join-the-mainnet)
  - [Join a testnet](#join-a-testnet)
  - [Run a local testnet](#run-a-local-testnet)
  - [Run a single node testnet](#run-a-single-node-testnet)
- [Set up a production environment](#set-up-a-production-environment)
  - [Increase maximum open files](#increase-maximum-open-files)
  - [Create a dedicated user](#create-a-dedicated-user)
  - [Port configuration](#port-configuration)
  - [Run the server as a daemon](#run-the-server-as-a-daemon)
  - [Register iqd as a service](#register-iqd-as-a-service)
  - [Start, stop, or restart service](#start-stop-or-restart-service)
  - [Access logs](#access-logs)
- [Resources](#resources)
- [Community](#community)
- [Contributing](#contributing)
- [License](#license)

## What is IQ Chain?

**[IQ Chain](https://iqchain.network)** is a public, open-source blockchain protocol that provides fundamental infrastructure for a decentralized economy and enables open participation in the creation of new financial primitives to power the innovation of money.

The IQ blockchain is secured by distributed consensus on staked asset BIQ and natively supports the issuance of [price-tracking stablecoins](https://docs.iqchain.network/docs/learn/glossary.html#algorithmic-stablecoin) that are algorithmically pegged to major world currencies, such as UST, KRT, and SDT. Smart contracts on IQChain run on WebAssembly and take advantage of core modules, such as on-chain swaps, price oracle, and staking rewards, to power modern [DeFi](https://docs.iqchain.network/docs/learn/glossary.html#defi) apps. Through dynamic fiscal policy managed by community governance, IQ is an evolving, democratized economy directed by its users.

**IQ Core** is the reference implementation of the IQ protocol, written in Golang. IQ Core is built atop [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and uses [Tendermint](https://github.com/tendermint/tendermint) BFT consensus. If you intend to work on IQ Core source, it is recommended that you familiarize yourself with the concepts in those projects.

## Installation

### Binaries

The easiest way to get started is by downloading a pre-built binary for your operating system. You can find the latest binaries on the [releases](https://github.com/bitwebs/core/releases) page.

### From Source

**Step 1. Install Golang**

Go v1.17+ or higher is required for IQ Core.

If you haven't already, install Golang by following the [official docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

**Step 2: Get IQ Core source code**

Use `git` to retrieve IQ Core from the [official repo](https://github.com/bitwebs/core/) and checkout the `main` branch. This branch contains the latest stable release, which will install the `iqd` binary.

```bash
git clone https://github.com/bitwebs/core/
cd core
git checkout main
```

**Step 3: Build IQ core**

Run the following command to install the executable `iqd` to your `GOPATH` and build IQ Core. `iqd` is the node daemon and CLI for interacting with a IQ node.

```bash
# COSMOS_BUILD_OPTIONS=rocksdb make install
make install
```

**Step 4: Verify your installation**

Verify that you've installed iqd successfully by running the following command:

```bash
iqd version --long
```

If iqd is installed correctly, the following information is returned:

```bash
name: iq
server_name: iqd
version: 0.5.0-rc0-9-g640fd0ed
commit: 640fd0ed921d029f4d1c3d88435bd5dbd67d14cd
build_tags: netgo,ledger
go: go version go1.17.2 darwin/amd64
```

## `iqd`

**NOTE:** `iqcli` has been deprecated and all of its functionalities have been merged into `iqd`.

`iqd` is the all-in-one command for operating and interacting with a running IQ node. For comprehensive coverage on each of the available functions, see [the iqd reference information](https://docs.iqchain.network/docs/develop/how-to/iqd/README.html). To view various subcommands and their expected arguments, use the `$ iqd --help` command:

<pre>
        <div align="left">
        <b>$ iqd --help</b>

        Stargate IQ App

        Usage:
          iqd [command]

        Available Commands:
          add-genesis-account Add a genesis account to genesis.json
          collect-gentxs      Collect genesis txs and output a genesis.json file
          debug               Tool for helping with debugging your application
          export              Export state to JSON
          gentx               Generate a genesis tx carrying a self delegation
          help                Help about any command
          init                Initialize private validator, p2p, genesis, and application configuration files
          keys                Manage your application's keys
          migrate             Migrate genesis to a specified target version
          query               Querying subcommands
          rosetta             spin up a rosetta server
          start               Run the full node
          status              Query remote node for status
          tendermint          Tendermint subcommands
          testnet             Initialize files for a iqd testnet
          tx                  Transactions subcommands
          unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
          validate-genesis    validates the genesis file at the default location or at the location passed as an arg
          version             Print the application binary version information

        Flags:
          -h, --help                help for iqd
              --home string         directory for config and data (default "/Users/$HOME/.iq")
              --log_format string   The logging format (json|plain) (default "plain")
              --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
              --trace               print out full stack trace on errors

        <b>Use "iqd [command] --help" for more information about a command.</b>
        </div>
</pre>

## Node Setup

Once you have `iqd` installed, you will need to set up your node to be part of the network.

### Join the mainnet

The following requirements are recommended for running a `columbus-5` mainnet node:

- **4 or more** CPU cores
- At least **2TB** of disk storage
- At least **100mbps** network bandwidth
- An Linux distribution

For configuration and migration instructions for setting up a Columbus-5 mainnet node, visit [The mainnet repo](https://github.com/bitwebs/mainnet).

**IQ Node Quick Start**
```
iqd init nodename
wget -O ~/.iq/config/genesis.json https://cloudflare-ipfs.com/ipfs/QmZAMcdu85Qr8saFuNpL9VaxVqqLGWNAs72RVFhchL9jWs
curl https://iqchain.network/addrbook.json > ~/.iqd/config/addrbook.json
iqd start
```

### Join a testnet

Several testnets might exist simultaneously. Ensure that your version of `iqd` is compatible with the network you want to join.

To set up a node on the latest testnet, visit [the testnet repo](https://github.com/bitwebs/testnet).

### Run a local testnet

The easiest way to set up a local testing environment is to run [LocalIQ](https://github.com/bitwebs/LocalIQ), which automatically orchestrates a complete testing environment suited for development with zero configuration.

### Run a single node testnet

You can also run a local testnet using a single node. On a local testnet, you will be the sole validator signing blocks.


**Step 1. Create network and account**

First, initialize your genesis file to bootstrap your network. Create a name for your local testnet and provide a moniker to refer to your node:

```bash
iqd init --chain-id=<testnet_name> <node_moniker>
```

Next, create a IQ account by running the following command:

```bash
iqd keys add <account_name>
```

**Step 2. Add account to genesis**

Next, add your account to genesis and set an initial balance to start. Run the following commands to add your account and set the initial balance:

```bash
iqd add-genesis-account $(iqd keys show <account_name> -a) 100000000ubiq,1000busd
iqd gentx <account_name> 10000000ubiq --chain-id=<testnet_name>
iqd collect-gentxs
```

**Step 3. Run IQ daemon**

Now you can start your private IQ network:

```bash
iqd start
```

Your `iqd` node will be running a node on `tcp://localhost:26656`, listening for incoming transactions and signing blocks.

Congratulations, you've successfully set up your local IQ network!

## Set up a production environment

**NOTE**: This guide only covers general settings for a production-level full node. You can find further details on considerations for operating a validator node by visiting the [IQ validator guide](https://docs.iqchain.network/docs/full-node/manage-a-iq-validator/README.html).

This guide has been tested against Linux distributions only. To ensure you successfully set up your production environment, consider setting it up on an Linux system.

### Increase maximum open files

`iqd` can't open more than 1024 files (the default maximum) concurrently.

You can increase this limit by modifying `/etc/security/limits.conf` and raising the `nofile` capability.

```
*                soft    nofile          65535
*                hard    nofile          65535
```

### Create a dedicated user

It is recommended that you run `iqd` as a normal user. Super-user accounts are only recommended during setup to create and modify files.

### Port configuration

`iqd` uses several TCP ports for different purposes.

- `26656`: The default port for the P2P protocol. Use this port to communicate with other nodes. While this port must be open to join a network, it does not have to be open to the public. Validator nodes should configure `persistent_peers` and close this port to the public.

- `26657`: The default port for the RPC protocol. This port is used for querying / sending transactions and must be open to serve queries from `iqd`. **DO NOT** open this port to the public unless you are planning to run a public node.

- `1317`: The default port for [Lite Client Daemon](https://docs.iqchain.network/docs/develop/how-to/start-lcd.html) (LCD), which can be enabled in `~/.iq/config/app.toml`. The LCD provides an HTTP RESTful API layer to allow applications and services to interact with your `iqd` instance through RPC. Check the [IQ REST API](https:/api.iqchain.network/swagger/#/) for usage examples. Don't open this port unless you need to use the LCD.

- `26660`: The default port for interacting with the [Prometheus](https://prometheus.io) database. You can use Promethues to monitor an environment. This port is closed by default.

### Run the server as a daemon

**Important**:

Keep `iqd` running at all times. The simplest solution is to register `iqd` as a `systemd` service so that it automatically starts after system reboots and other events.


### Register iqd as a service

First, create a service definition file in `/etc/systemd/system`.

**Sample file: `/etc/systemd/system/iqd.service`**

```
[Unit]
Description=IQ Daemon
After=network.target

[Service]
Type=simple
User=iq
ExecStart=/data/iq/go/bin/iqd start
Restart=on-abort

[Install]
WantedBy=multi-user.target

[Service]
LimitNOFILE=65535
```

Modify the `Service` section from the given sample above to suit your settings.
Note that even if you raised the number of open files for a process, you still need to include `LimitNOFILE`.

After creating a service definition file, you should execute `systemctl daemon-reload`.

### Start, stop, or restart service

Use `systemctl` to control (start, stop, restart)

```bash
# Start
systemctl start iqd
# Stop
systemctl stop iqd
# Restart
systemctl restart iqd
```

### Access logs

```bash
# Entire log
journalctl -t iqd
# Entire log reversed
journalctl -t iqd -r
# Latest and continuous
journalctl -t iqd -f
```

## Resources

- Developer Tools

  - IQ developer documentation(https://docs.iqchain.network)
  - SDKs
    - [IQ.js](https://www.github.com/bitwebs/iq.js) for JavaScript
    - [iq-sdk-python](https://www.github.com/bitwebs/iq-sdk-python) for Python
  - [Faucet](https://faucet.iqchain.network) can be used to get tokens for testnets
  - [LocalIQ](https://www.github.com/bitwebs/LocalIQ) can be used to set up a private local testnet with configurable world state

- Block Explorers

  - [IQ Finder](https://finder.iqchain.network) - IQ's basic block explorer.

- Wallets

  - [IQ Station](https://station.iqchain.network) - The official IQ wallet.
  - IQ Station Mobile
    - [iOS](https://apps.apple.com/us/app)
    - [Android](https://play.google.com/store/apps/details)

## Community

- [Offical Website](https://iqchain.network)
- [Telegram](https://t.me/iq_announcements)
- [Twitter](https://twitter.com/iqchain)
- [YouTube](https://goo.gl/3G4T1z)

## Contributing

If you are interested in contributing to IQ Core source, please review our [code of conduct](./CODE_OF_CONDUCT.md).

## License

This software is licensed under the Apache 2.0 license. Read more about it [here](LICENSE).

© 2021 BitWeb Labs, PTE LTD

<hr/>

<p>&nbsp;</p>
<div align="center">
  <sub><em>Powering the innovation of money.</em></sub>
</div>
