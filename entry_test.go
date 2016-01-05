package factom_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	ed "github.com/agl/ed25519"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprint("testing")

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

func TestComposeEntryCommit(t *testing.T) {
	pub, pri, err := ed.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	
	e := factom.NewEntry()
	if err := e.UnmarshalJSON(jsonentry); err != nil {
		t.Error(err)
	}
	j, err := factom.ComposeEntryCommit(pub, pri, e)
	if err != nil {
		t.Error(err)
	}
	
	t.Log("json:", string(j))
}

func TestComposeEntryReveal(t *testing.T) {
	e := factom.NewEntry()
	if err := e.UnmarshalJSON(jsonentry); err != nil {
		t.Error(err)
	}
	
	j, err := factom.ComposeEntryReveal(e)
	if err != nil {
		t.Error(err)
	}
	
	t.Log("json:", string(j))
}

func TestComposeChainCommit(t *testing.T) {
	pub, pri, err := ed.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	
	e := factom.NewEntry()
	if err := e.UnmarshalJSON(jsonentry); err != nil {
		t.Error(err)
	}
	
	c := factom.NewChain(e)
	
	j, err := factom.ComposeChainCommit(pub, pri, c)
	if err != nil {
		t.Error(err)
	}
	
	t.Log("json:", string(j))
}
