// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

type addressRequest struct {
	Address string `json:"address"`
}

type importRequest struct {
	Addresses []struct {
		Secret string `json:secret`
	} `json:addresses`
}

type addressResponse struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

type multiAddressResponse struct {
	Addresses []*addressResponse `json:"addresses"`
}

type walletBackupResponse struct {
	Seed      string             `json:wallet-seed`
	Addresses []*addressResponse `json:"addresses"`
}
