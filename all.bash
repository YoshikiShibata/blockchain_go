#!/bin/bash -x 

# Copyright (C) 2017 Yoshiki Shibata. All rights reserved.

echo "Removing the blockchain db ..."
rm blockchain.db
echo ""

echo "Building blockchain program as bc ..."
go build -o bc
echo ""

echo "Create the blockchain for Ivan ..."
time ./bc createblockchain -address Ivan
echo ""

echo "Get the balance for Ivan ..."
./bc getbalance -address Ivan
echo ""

echo "Send 6 amount of money from Iva to Pedro"
time ./bc send -from Ivan -to Pedro -amount 6
echo ""

echo "Get the balance for Ivan ..."
./bc getbalance -address Ivan
echo ""

echo "Get the balance for Pedro ..."
./bc getbalance -address Pedro
echo ""

