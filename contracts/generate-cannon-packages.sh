#!/bin/bash

set -e

network=$1
chainId=$2
copyFrom=$3

CANNON=${CANNON:-cannon}

feeddata=$(curl -q https://reference-data-directory.vercel.app/feeds-${network}.json)

for feedIdx in $(echo "$feeddata" | jq -r 'keys | join(" ")'); do
    feedname=$(echo "$feeddata" | jq -r ".[${feedIdx}].ens" | tr -d '-' | tr '[:lower:]' '[:upper:]')
    feedproxy=$(echo "$feeddata" | jq -r ".[${feedIdx}].proxyAddress")
    feedimpl=$(echo "$feeddata" | jq -r ".[${feedIdx}].contractAddress")

    echo "configuring package $feedname ($feedproxy)"

    $CANNON alter chainlink-aggregator:1.0.0 --chain-id $chainId --preset $feedname set-url $copyFrom
    $CANNON alter chainlink-aggregator:1.0.0 --chain-id $chainId --preset $feedname set-contract-address Proxy $feedproxy
    $CANNON alter chainlink-aggregator:1.0.0 --chain-id $chainId --preset $feedname set-contract-address AggregatorImpl $feedimpl
done
