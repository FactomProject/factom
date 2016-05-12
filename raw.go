// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/hex"
)

type RawData struct {
	Data string `json:"data"`
}

func (r *RawData) GetDataBytes() ([]byte, error) {
	return hex.DecodeString(r.Data)
}
