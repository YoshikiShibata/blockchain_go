#!/bin/bash

# Copyright (C) 2017 Yoshiki Shibata. All rights reserved.

echo "Removing the blockchain db ..."
rm blockchain.db
echo ""

echo "Building blockchain program as bc ..."
go build -o bc

if [ $? != 0 ];
then
	exit 1
fi

echo "Creating a wallet ..."
W1=`./bc createwallet`
echo ${W1}
echo ""

echo "Creating a wallet ..."
W2=`./bc createwallet`
echo ${W2}
echo ""

echo "Creating a wallet ..."
W3=`./bc createwallet`
echo ${W3}
echo ""

ADD1=`echo ${W1} | awk '{ print $4 }'`
ADD2=`echo ${W2} | awk '{ print $4 }'`
ADD3=`echo ${W3} | awk '{ print $4 }'`

echo "Create a blockchain ..."
./bc createblockchain -address ${ADD1}
echo ""

echo "Get balace for ${ADD1} .."
./bc getbalance -address ${ADD1}
echo ""

echo "Send amount 5 of money from ${ADD2} to ${ADD1} ..."
./bc send -from ${ADD2} -to ${ADD1} -amount 5
echo ""

echo "Send amount 6 of money from ${ADD1} to ${ADD2} ..."
./bc send -from ${ADD1} -to ${ADD2} -amount 6
echo ""

echo "Get Balance for ${ADD1} ..."
./bc getbalance -address ${ADD1}
echo ""

echo "Get Balance for ${ADD2} ..."
./bc getbalance -address ${ADD2}
echo ""

echo "Get Balance for ${ADD3} ..."
./bc getbalance -address ${ADD3}
echo ""

