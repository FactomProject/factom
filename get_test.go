// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"
	"math/rand"
	"testing"
)

var _ = fmt.Sprint("testing")

const gutenbergChainID = "00511c298668bc5032a64b76f8ede6f119add1a64482c8602966152c0b936c77"

func TestGetAllChainEntries(t *testing.T) {
	t.Skip("Skip this test in short mode")
	es, err := GetAllChainEntries(gutenbergChainID)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(es))
	t.Log(es[rand.Intn(len(es))])
	t.Log(es[rand.Intn(len(es))])
}

func TestGetFirstEntry(t *testing.T) {
	e, err := GetFirstEntry(gutenbergChainID)
	if err != nil {
		t.Error(err)
	}
	t.Log(e)
}
