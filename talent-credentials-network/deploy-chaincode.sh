#!/bin/bash

export CC_PACKAGE_SEQUENCE=${1:-1}

# Check if all required scripts are available
if [ ! -f ./create-chaincode-pkg.sh ]; then
    echo "Error: create-chaincode-pkg.sh not found!"
    exit 1
fi

if [ ! -f ./install-chaincode.sh ]; then
    echo "Error: install-chaincode.sh not found!"
    exit 1
fi

if [ ! -f ./approve-chaincode.sh ]; then
    echo "Error: approve-chaincode.sh not found!"
    exit 1
fi

if [ ! -f ./commit-chaincode.sh ]; then
    echo "Error: commit-chaincode.sh not found!"
    exit 1
fi

if [ ! -f ./invoke-chaincode.sh ]; then
    echo "Error: invoke-chaincode.sh not found!"
    exit 1
fi

echo "Starting chaincode deployment process..."

# Step 1: Create Chaincode Package
echo "Step 1: Creating chaincode package..."
./create-chaincode-pkg.sh
if [ $? -ne 0 ]; then
    echo "Error: Chaincode package creation failed."
    exit 1
fi

# Step 2: Install Chaincode on peers
echo "Step 2: Installing chaincode on peers..."
./install-chaincode.sh
if [ $? -ne 0 ]; then
    echo "Error: Chaincode installation failed."
    exit 1
fi

# Step 3: Approve Chaincode for both Orgs
echo "Step 3: Approving chaincode for both Orgs..."
./approve-chaincode.sh $CC_PACKAGE_SEQUENCE
if [ $? -ne 0 ]; then
    echo "Error: Chaincode approval failed."
    exit 1
fi

# Step 5: Commit Chaincode to the Channel
echo "Step 5: Committing chaincode to the channel..."
./commit-chaincode.sh $CC_PACKAGE_SEQUENCE
if [ $? -ne 0 ]; then
    echo "Error: Chaincode commit failed."
    exit 1
fi

# Step 6: Invoke Chaincode to initialize ledger
echo "Step 6: Invoking chaincode to initialize the ledger..."
./invoke-chaincode.sh
if [ $? -ne 0 ]; then
    echo "Error: Chaincode invocation failed."
    exit 1
fi

echo "Chaincode deployment completed successfully!"
