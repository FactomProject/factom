package factom_test

import (
	"testing"
	
	. "github.com/FactomProject/factom"
)

func TestResolveDnsName(t *testing.T) {
	f1, e1, err1 := ResolveDnsName(goodAddr)
	if err1 != nil {
		t.Error(err1)
	}
	t.Logf("fct: %s\nec: %s\n", f1, e1)
	
	f2, e2, err2 := ResolveDnsName(badAddr)
	if err2 == nil {
		t.Errorf("bad address %s did not return error", badAddr)
	}
	t.Logf("fct: %d\nec: %d\n", f2, e2)
}