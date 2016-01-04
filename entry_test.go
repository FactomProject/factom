package factom_test

import (
	"testing"
	"github.com/FactomProject/factom"
)

var jsonentry = []byte(`
{
	"ChainID":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	"ExtIDs":[
		"foo",
		"bar"
	],
	"Content":"Hello Factom!"
}`)

func TestUnmarshalJSON(t *testing.T) {
	e := factom.NewEntry()
	if err := e.UnmarshalJSON(jsonentry); err != nil {
		t.Error(err)
	}
	t.Log(e)
}