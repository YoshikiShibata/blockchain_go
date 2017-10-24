#!/bin/bash

# Copyright (C) 2017 Yoshiki Shibata. All rights reserved.

export NODE_ID=3000

echo "Running Central Node"
echo ""

rm wallet_${NODE_ID}.dat blockchain_${NODE_ID}.db

echo "Creating a wallet ..."
W=`./bc createwallet`
echo ${W}
echo ""

ADD=`echo ${W} | awk '{ print $4 }'`

echo "Create a blockchain ..."
./bc createblockchain -address ${ADD}
echo ""

echo "Copy the db as a Genesis db"
cp blockchain_${NODE_ID}.db blockchain_genesis.db

