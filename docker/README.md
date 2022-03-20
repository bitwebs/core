# IQ Core Dockerized

## Common usage examples 

### Mainnet 

Standard options with LCD enabled: 

```
docker run -it -p 1317:1317 -p 26657:26657 -p 26656:26656 bitwebs/iq-core-node:v0.0.1-boom
```

LCD disabled: 

```
docker run -e ENABLE_LCD=false -it -p 1317:1317 -p 26657:26657 -p 26656:26656 bitwebs/iq-core-node:v0.0.1-boom
```

Custom gas fees: 

```
docker run -e MINIMUM_GAS_PRICES="0.01133ubiq,0.15ubusd,0.104938ubsdr,169.77ubkrw,428.571ubmnt,0.125ubeur,0.98ubcny,16.37ubjpy,0.11ubgbp,10.88ubinr,0.19ubcad,0.14ubchf,0.19ubaud,0.2ubsgd,4.62ubthb,1.25ubsek,1.25ubnok,0.9ubdkk,2180.0ubidr,7.6ubphp,1.17ubhkd" -it -p 1317:1317 -p 26657:26657 -p 26656:26656 bitwebs/iq-core-node:v0.0.1-boom
```

Starting the sync from a snapshot:

```
docker run -e SNAPSHOT_NAME="swartz-1.tar.lz4" -it -p 1317:1317 -p 26657:26657 -p 26656:26656 bitwebs/iq-core-node:v0.0.1-boom
```

You can find the latest snapshots [here](https://quicksync.io/networks/terra.html).

Custom snapshot URL:

```
docker run -e SNAPSHOT_BASE_URL="https://get.quicksync.io" -it -p 1317:1317 -p 26657:26657 -p 26656:26656 bitwebs/iq-core-node:v0.0.1-boom
```

**Note:** We recommend copying a snapshot to S3 or another file store and using the above options to point the container to your snapshot. The default snapshot name included will be obsolete and removed in a matter of days.

Starting a McAfee node: 

```
docker run -it -p 1317:1317 -p 26657:26657 -p 26656:26656 bitwebs/iq-core-node:v0.0.1-boom-testnet
```

## Building the Docker images

```
./build_all.sh v0.0.1-boom
```