// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"

	netki "github.com/FactomProject/netki-go-partner-client"
)

<<<<<<< HEAD
func ResolveDnsName(addr string) (fct, ec string, err error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/resolve-address/%s",
			serverFct,
			addr))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	type x struct {
		Response string
		Success  bool
	}
	a := new(x)
	if err = json.Unmarshal(body, a); err != nil {
		return
	}
	if !a.Success {
		err = fmt.Errorf(a.Response)
		return
=======
func ResolveDnsName(addr string) (string, string, error) {
	fct, err1 := netki.WalletNameLookup(addr, "fct")
	ec, err2 := netki.WalletNameLookup(addr, "fec")
	if err1 != nil && err2 != nil {
		return fct, ec, fmt.Errorf("%s\n%s", err1, err2)
	}
	return fct, ec, nil
}

func GetDnsBalance(addr string) (int64, int64, error) {
	fct, ec, err := ResolveDnsName(addr)
	if err != nil {
		return -1, -1, err
>>>>>>> FactomProject/m2-merge
	}

	f, err1 := GetFactoidBalance(fct)
	e, err2 := GetECBalance(ec)
	if err1 != nil || err2 != nil {
		return f, e, fmt.Errorf("%s\n%s\n", err1, err2)
	}
<<<<<<< HEAD
	b := new(y)
	if err = json.Unmarshal([]byte(a.Response), b); err != nil {
		return
	}

	return b.Fct, b.Ec, nil
=======

	return f, e, nil
>>>>>>> FactomProject/m2-merge
}
