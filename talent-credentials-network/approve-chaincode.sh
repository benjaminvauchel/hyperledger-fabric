#!/bin/bash

export CC_PACKAGE_VERSION=1.5
export CC_PACKAGE_SEQUENCE=${1:-1}

setEnvOrg1() {
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_ADDRESS=localhost:7051
}

setEnvOrg2() {
  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
  export CORE_PEER_ADDRESS=localhost:9051
}

echo "Querying installed chaincodes to get the package ID..."
export CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "basic_1.5" | awk -F"Package ID: " '{print $2}' | awk '{print $1}' | sed 's/,//g')

echo "Package ID: $CC_PACKAGE_ID"

# Approve chaincode definition for Org1
setEnvOrg1
echo "Approving chaincode definition for Org1..."
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version $CC_PACKAGE_VERSION --package-id $CC_PACKAGE_ID --sequence $CC_PACKAGE_SEQUENCE --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# Approve chaincode definition for Org2
setEnvOrg2
echo "Approving chaincode definition for Org2..."
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version $CC_PACKAGE_VERSION --package-id $CC_PACKAGE_ID --sequence $CC_PACKAGE_SEQUENCE --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

echo "Chaincode definition approved for both Org1 and Org2"
