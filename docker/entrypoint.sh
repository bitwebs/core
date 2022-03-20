#!/bin/sh

# Default to "data".
DATADIR="${DATADIR:-/root/.iq/data}"
MONIKER="${MONIKER:-docker-node}"
ENABLE_LCD="${ENABLE_LCD:-true}"
MINIMUM_GAS_PRICES=${MINIMUM_GAS_PRICES-0.01133ubiq,0.15ubusd,0.104938ubsdr,169.77ubkrw,428.571ubmnt,0.125ubeur,0.98ubcny,16.37ubjpy,0.11ubgbp,10.88ubinr,0.19ubcad,0.14ubchf,0.19ubaud,0.2ubsgd,4.62ubthb,1.25ubsek,1.25ubnok,0.9ubdkk,2180.0ubidr,7.6ubphp,1.17ubhkd}
SNAPSHOT_NAME="${SNAPSHOT_NAME}"
SNAPSHOT_BASE_URL="${SNAPSHOT_BASE_URL:-https://getsfo.quicksync.io}"

# First sed gets the app.toml moved into place.
# app.toml updates
sed 's/minimum-gas-prices = "0ubiq"/minimum-gas-prices = "'"$MINIMUM_GAS_PRICES"'"/g' ~/app.toml > ~/.iq/config/app.toml

# Needed to use awk to replace this multiline string.
if [ "$ENABLE_LCD" = true ] ; then
  gawk -i inplace '/^# Enable defines if the API server should be enabled./,/^enable = false/{if (/^enable = false/) print "# Enable defines if the API server should be enabled.\nenable = true"; next} 1' ~/.iq/config/app.toml
fi

# config.toml updates

sed 's/moniker = "moniker"/moniker = "'"$MONIKER"'"/g' ~/config.toml > ~/.iq/config/config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' ~/.iq/config/config.toml

if [ "$CHAINID" = "swartz-1" ] && [[ ! -z "$SNAPSHOT_NAME" ]] ; then 
  # Download the snapshot if data directory is empty.
  res=$(find "$DATADIR" -name "*.db")
  if [ "$res" ]; then
      echo "data directory is NOT empty, skipping quicksync"
  else
      echo "starting snapshot download"
      mkdir -p $DATADIR
      cd $DATADIR
      FILENAME="$SNAPSHOT_NAME"

      # Download
      aria2c -x5 $SNAPSHOT_BASE_URL/$FILENAME
      # Extract
      lz4 -d $FILENAME | tar xf -

      # # cleanup
      rm $FILENAME
  fi
fi

exec "$@" --db_dir $DATADIR