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

echo "Sending 2 amount of money from Pedro Helen"
./bc send -from Pedro -to Helen -amount 2
echo ""

echo "Sending 2 amount of moeny from Ivan to Helen"
./bc send -from Ivean -to Helen -amount 2

echo "Sending 3 amount of moeney from Helen to Rachel"
./bc send -from Helen to Rachel -amount 2

./bc getbalance -address Ivan
./bc getbalance -address Pedro
./bc getbalance -address Helen
./bc getbalance -address Rachel

./bc send -from Pedro -to Ivan -amount 5
./bc getbalance -address Pedro
./bc getbalance -address Ivan
