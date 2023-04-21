#!/bin/bash

set -e

network=$1
chainId=$2
copyFrom=$3

CANNON=${CANNON:-cannon}

# chainlink token deployment is a dependency of the price feeds
echo "Configuring chainlink token"
$CANNON alter chainlink-token:1.0.0 --chain-id $chainId set-url ipfs://QmTogXYidVGmJSxR3fg3kHsDAT8eYvmDp9y3zpu8T2VXub

if [ "$network" = "ethereum-testnet-sepolia" ]; then
    $CANNON alter chainlink-token:1.0.0 --chain-id $chainId set-contract-address Token '0x779877A7B0D9E8603169DdbD7836e478b4624789'
elif [ "$network" = "mainnet" ]; then
    $CANNON alter chainlink-token:1.0.0 --chain-id $chainId set-contract-address Token '0x514910771AF9Ca656af840dff83E8264EcF986CA'
fi

feeddata=$(curl -q https://reference-data-directory.vercel.app/feeds-${network}.json)

set +e
for feedIdx in $(echo "$feeddata" | jq -r 'keys | join(" ")'); do
    feedname=$(echo "$feeddata" | jq -r ".[${feedIdx}].ens" | tr -d '-' | tr '[:upper:]' '[:lower:]')
    feedproxy=$(echo "$feeddata" | jq -r ".[${feedIdx}].proxyAddress")
    feedimpl=$(echo "$feeddata" | jq -r ".[${feedIdx}].contractAddress")

    echo "configuring package $feedname ($feedproxy)"

    $CANNON alter chainlink-aggregator:1.0.0 --chain-id $chainId --preset $feedname set-url $copyFrom
    $CANNON alter chainlink-aggregator:1.0.0 --chain-id $chainId --preset $feedname set-contract-address Proxy $feedproxy
    $CANNON alter chainlink-aggregator:1.0.0 --chain-id $chainId --preset $feedname set-contract-address AggregatorImpl $feedimpl
done
set -e

# chainlink VRF depends on the price feeds

echo "Configuring chainlink VRF"
$CANNON alter chainlink-vrf:2.0.0 --chain-id $chainId set-url ipfs://QmWkwqf6Vd6x17A32W5iQQQBMMSnhqw1RbZvbM8dzFaxG6

if [ "$network" = "ethereum-testnet-sepolia" ]; then
    $CANNON alter chainlink-vrf:2.0.0 --chain-id $chainId set-contract-address VRFCoordinator '	0x8103B0A8A00be2DDC778e6e7eaa21791Cd364625'
    $CANNON alter chainlink-vrf:2.0.0 --chain-id $chainId set-contract-address VRFWrapper '0xab18414CD93297B0d12ac29E63Ca20f515b3DB46'
elif [ "$network" = "mainnet" ]; then
    $CANNON alter chainlink-vrf:2.0.0 --chain-id $chainId set-contract-address VRFCoordinator '0x271682DEB8C4E0901D1a1550aD2e64D568E69909'
    $CANNON alter chainlink-vrf:2.0.0 --chain-id $chainId set-contract-address VRFWrapper '0x5A861794B927983406fCE1D062e00b9368d97Df6'
fi
