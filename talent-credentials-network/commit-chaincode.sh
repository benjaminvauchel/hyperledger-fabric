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

# Query installed chaincode to get the package ID
echo "Querying installed chaincodes to get the package ID..."
export CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep "basic_1.5" | awk -F"Package ID: " '{print $2}' | awk '{print $1}' | sed 's/,//g')

echo "Package ID: $CC_PACKAGE_ID"

# Verify the commit readiness
echo "Checking commit readiness..."
setEnvOrg1
commit_ready=$(peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version $CC_PACKAGE_VERSION --sequence $CC_PACKAGE_SEQUENCE --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json)

# Print the full commit readiness response to debug
echo "Commit readiness response:"
echo $commit_ready

# Check if all organizations have approved the chaincode definition
approved_org1=$(echo $commit_ready | jq -r '.approvals.Org1MSP')
approved_org2=$(echo $commit_ready | jq -r '.approvals.Org2MSP')

# Trim any extra spaces from the results
approved_org1=$(echo $approved_org1 | xargs)
approved_org2=$(echo $approved_org2 | xargs)

echo "Org1MSP approval: $approved_org1"
echo "Org2MSP approval: $approved_org2"

if [[ "$approved_org1" == "true" && "$approved_org2" == "true" ]]; then
  echo "Both organizations have approved the chaincode definition."

  # Commit the chaincode
  echo "Committing chaincode definition to the channel..."

  # Approve chaincode for Org1
  setEnvOrg1
  peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version $CC_PACKAGE_VERSION --sequence $CC_PACKAGE_SEQUENCE --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

  echo "Chaincode committed to the channel."

  # Query committed chaincode
  echo "Querying the committed chaincode..."
  peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

else
  echo "The chaincode definition has not been approved by both organizations. Please ensure both organizations approve the chaincode definition."
  exit 1
fi
