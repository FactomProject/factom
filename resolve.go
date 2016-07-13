// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"
	
	netki "github.com/FactomProject/netki-go-partner-client"
)

func ResolveDnsName(addr string) (string, string, error) {
	fct, ferr := netki.WalletNameLookup(addr, "fct")
	ec, eerr := netki.WalletNameLookup(addr, "fec")
	if ferr != nil && eerr != nil {
		return fct, ec, fmt.Errorf("%s\n%s", ferr, eerr)
	}
	return fct, ec, nil
}

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
