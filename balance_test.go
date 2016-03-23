package factom_test

import (
	"fmt"
	"testing"

	. "github.com/FactomProject/factom"
)

var (
	s        = fmt.Sprint("testing")
	badAddr  = "bad.factom.bit"
	goodAddr = "factom.michaeljbeam.me"
)

func TestDnsBalance(t *testing.T) {
	f1, e1, err1 := DnsBalance(badAddr)
	t.Logf("fct: %d\nec: %d\n", f1, e1)
	if err1 == nil {
		t.Errorf("bad address %s did not return error", badAddr)
	}

	f2, e2, err2 := DnsBalance(goodAddr)
	t.Logf("fct: %d\nec: %d\n", f2, e2)
	if err2 != nil {
		t.Error(err2)
	}
}
