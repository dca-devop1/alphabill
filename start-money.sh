#!/bin/bash
# build binary
make clean build
mkdir testab
mkdir testab/rootchain
nodeAddresses=""

# generate fee bill for money-genesis
feeBill='{"systemId": "0x00000000", "unitId": "0x0000000000000000000000000000000000000000000000000000000000000002", "ownerPubKey": "0x03c30573dc0c7fd43fcb801289a6a96cb78c27f4ba398b89da91ece23e9a99aca3"}'
echo $feeBill > testab/money-fee-bill.json

# Generate node genesis files.
for i in 1 2 3
do
  # "-g" flags also generates keys
  build/alphabill money-genesis --home testab/money$i -g -c testab/money-fee-bill.json
done

# generate rootchain and partition genesis files
build/alphabill root-genesis --home testab/rootchain -o testab/rootchain/genesis -p testab/money1/money/node-genesis.json -p testab/money2/money/node-genesis.json -p testab/money3/money/node-genesis.json -k testab/rootchain/keys.json -g

#start root chain
build/alphabill root --home testab/rootchain -f testab/rootchain/rounds.db -k testab/rootchain/keys.json -g testab/rootchain/genesis/root-genesis.json > testab/rootchain/rootchain.log &

port=26666
# partition node addresses
for i in 1 2 3
do
  id=$(build/alphabill identifier -k testab/money$i/money/keys.json | tail -n1)
  nodeAddresses="$nodeAddresses,$id=/ip4/127.0.0.1/tcp/$port";

  ((port=port+1))
done

nodeAddresses="${nodeAddresses:1}"

port=26666
grpcPort=26766
restPort=26866
#start partition nodes
for i in 1 2 3
do
  build/alphabill money --home testab/money$i -f testab/money$i/money/blocks.db -k testab/money$i/money/keys.json -r "/ip4/127.0.0.1/tcp/26662" -a "/ip4/127.0.0.1/tcp/$port" --server-address ":$grpcPort" --rest-server-address "localhost:$restPort" -g testab/rootchain/genesis/partition-genesis-0.json -p "$nodeAddresses" > "testab/money$i/money$i.log" &
  ((port=port+1))
  ((grpcPort=grpcPort+1))
  ((restPort=restPort+1))
done
