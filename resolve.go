// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"

	netki "github.com/FactomProject/netki-go-partner-client"
)

// ResolveDnsName resolve a netki wallet DNS name and returns the public address
// string for the Factoid address and the Entry Credit addresses associated with
// that name.
func ResolveDnsName(addr string) (string, string, error) {
	fct, err1 := netki.WalletNameLookup(addr, "fct")
	ec, err2 := netki.WalletNameLookup(addr, "fec")
	if err1 != nil && err2 != nil {
		return fct, ec, fmt.Errorf("%s\n%s", err1, err2)
	}
	return fct, ec, nil
}

// GetDnsBalance returns the balances of the Factoid and Entry Credit addresses
// associated with a netki DNS name.
func GetDnsBalance(addr string) (int64, int64, error) {
	fct, ec, err := ResolveDnsName(addr)
	if err != nil {
		return -1, -1, err
	}

	f, err1 := GetFactoidBalance(fct)
	e, err2 := GetECBalance(ec)
	if err1 != nil || err2 != nil {
		return f, e, fmt.Errorf("%s\n%s\n", err1, err2)
	}

	return f, e, nil
}
