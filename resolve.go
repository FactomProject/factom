// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	}

	type y struct {
		Fct, Ec string
	}
	b := new(y)
	if err = json.Unmarshal([]byte(a.Response), b); err != nil {
		return
	}
	
	return b.Fct, b.Ec, nil
}
