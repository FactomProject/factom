#!/bin/bash
#set -x

# This is to test the the load that wallet-balances can handle.

FCTADDWFACTOIDS="FA3EPZYqodgyEGXNMbiZKE5TS2x2J9wF8J9MvPZb52iGR78xMgCb"

# This takes makes 10,000 factoid addresses and funds all of them with a random factoshi balance 

COUNTFCT=10000
OUTPUTFCT=""

while [ $COUNT -gt 0 ]; do
	OUTPUTFCT="$(echo `factom-cli newfctaddress`)"
	AMOUNTFCT="1.$((RANDOM % 10000))"
	factom-cli sendfct FA3EPZYqodgyEGXNMbiZKE5TS2x2J9wF8J9MvPZb52iGR78xMgCb "$OUTPUTFCT" "$AMOUNTFCT" >> factomcli.log & 
	let COUNTFCT=COUNTFCT-1
done

# This creates 1000 entry credit addesses and funds them

COUNTEC=1000
OUTPUTEC=""

while [ $COUNTEC -gt 0 ]; do
	OUTPUTEC="$(echo `factom-cli newecaddress`)"
	AMOUNTEC="$((1 + RANDOM % 15))"
	factom-cli buyec "$FCTADDWFACTOIDS" "$OUTPUTEC" "$AMOUNTEC" >> factomcliec.log &
	let COUNTEC=COUNTEC-1
done
