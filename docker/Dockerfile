ARG version=v0.0.1-boom

FROM bitwebs/iq-core:${version}

ARG chainid=swartz-1

ENV CHAINID ${chainid}

# Moniker will be updated by entrypoint.
RUN iqd init --chain-id $chainid moniker

# Backup for templating
RUN mv ~/.iq/config/config.toml ~/config.toml
RUN mv ~/.iq/config/app.toml ~/app.toml

RUN if [ "$chainid" = "swartz-1" ] ; then wget -O ~/.iq/config/genesis.json https://iqchain.network/swartz-1-genesis.json; fi
RUN if [ "$chainid" = "swartz-1" ] ; then wget -O ~/.iq/config/addrbook.json https://iqchain.network/addrbook.json; fi

RUN if [ "$chainid" = "mcafee-1" ] ; then wget -O ~/.iq/config/genesis.json https://raw.githubusercontent.com/bitwebs/iq-testnets/master/mcafee-1/genesis.json; fi
RUN if [ "$chainid" = "mcafee-1" ] ; then wget -O ~/.iq/config/addrbook.json https://raw.githubusercontent.com/bitwebs/iq-testnets/master/mcafee-1/addrbook.json; fi

RUN apk update && apk add wget lz4 aria2 curl jq gawk coreutils

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

CMD ["/usr/local/bin/iqd", "start"]