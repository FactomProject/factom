// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/FactomProject/factom"
)

func TestUnmarshalFBlock(t *testing.T) {
	js := []byte(`{"fblock":{"bodymr":"0b6823522198d47689065e7b492baafbf817f0036934afffd1c968f2533a3e84","prevkeymr":"48c432b586b1737bc8ea0349ec319e41f07b28bc89d94b2e970e09f494eb8e04","prevledgerkeymr":"7a7c9851d9bcfb00f4d3d4cd0179adb43e47aabed628e7fceaf0ca718853045b","exchrate":90900,"dbheight":20002,"transactions":[{"txid":"fab98df81a80b1177c5226ff307be7ecc77c30666c63f06623a606424d41fe72","blockheight":0,"millitimestamp":1453149000985,"inputs":[],"outputs":[],"outecs":[],"rcds":[],"sigblocks":[]},{"txid":"1ec91421e01d95267f3deb9b9d5f29d3438387a0280a5ffa5e9a60f235212ae8","blockheight":0,"millitimestamp":1453149058599,"inputs":[{"amount":26268275436,"address":"3d956f129c08ac413025be3f6e47e3fb26461df35c9ccaf2fe4d53373e52536b","useraddress":"FA2SCdYb8iBYmMcmeUjHB8NhKx6DqH3wDovkumgbKt4oNkD3TJMg"}],"outputs":[{"amount":26267184636,"address":"ccf82cf94557f08a6859d8bf4a9b3ce361d0abae1e3bf5136b24638b74d32bc6","useraddress":"FA3XME5vdcjG8jPT188UFkum9BeAJJLgwyCkGB12QLsDA2qQaBET"}],"outecs":[],"rcds":["016664074524dd6a58e6593780717233b56d381a6798e5ee5ba75564bde589a6bf"],"sigblocks":[{"signatures":["efdab088b50d56ea2dfd4f600d5727a06cd7e9f3c353288e6898723ea32f4f044d27a80a199cfefec06cf53e18ea863b05b1075001d592b913e7f32c3d3f2204"]}]}],"chainid":"000000000000000000000000000000000000000000000000000000000000000f","keymr":"cfcac07b29ccfa413aeda646b5d386006468189939dfdfa6415b97cc35f2ea1a","ledgerkeymr":"a47da86f6ac8111da8a7d2a64fbaed1f74839722276acc5773b908963d01a029"},"rawdata":"000000000000000000000000000000000000000000000000000000000000000f0b6823522198d47689065e7b492baafbf817f0036934afffd1c968f2533a3e8448c432b586b1737bc8ea0349ec319e41f07b28bc89d94b2e970e09f494eb8e047a7c9851d9bcfb00f4d3d4cd0179adb43e47aabed628e7fceaf0ca718853045b000000000001631400004e220000000002000000c9020152566e1519000000020152566ef627010100e1edd8a56c3d956f129c08ac413025be3f6e47e3fb26461df35c9ccaf2fe4d53373e52536be1ed95db7cccf82cf94557f08a6859d8bf4a9b3ce361d0abae1e3bf5136b24638b74d32bc6016664074524dd6a58e6593780717233b56d381a6798e5ee5ba75564bde589a6bfefdab088b50d56ea2dfd4f600d5727a06cd7e9f3c353288e6898723ea32f4f044d27a80a199cfefec06cf53e18ea863b05b1075001d592b913e7f32c3d3f220400000000000000000000"}`)

	// Create temporary struct to unmarshal json object
	wrap := new(struct {
		FBlock  *FBlock `json:"fblock"`
		RawData []byte  `json:"rawdata"`
	})

	err := json.Unmarshal(js, wrap)
	if err != nil {
		t.Error(err)
	}
	t.Log(wrap.FBlock)
}

func TestGetFBlock(t *testing.T) {
	fb, raw, err := GetFBlock("cfcac07b29ccfa413aeda646b5d386006468189939dfdfa6415b97cc35f2ea1a")
	if err != nil {
		t.Error(err)
	}
	t.Log(fb)
	t.Log(fmt.Printf("%x\n", raw))
}
