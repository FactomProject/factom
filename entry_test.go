package factom_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprint("testing")

var jsonentry = []byte(`
{
	"ChainID":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	"ExtIDs":[
		"bbbb",
		"cccc"
	],
	"Content":"111111111111111111"
}`)

var jsonentry2 = []byte(`
{
	"ChainName":["aaaa", "bbbb"],
	"ExtIDs":[
		"cccc",
		"dddd"
	],
	"Content":"111111111111111111"
}`)

func TestUnmarshalJSON(t *testing.T) {
	e1 := factom.NewEntry()
	if err := e1.UnmarshalJSON(jsonentry); err != nil {
		t.Error(err)
	}
	t.Log(e1)

	e2 := factom.NewEntry()
	if err := e2.UnmarshalJSON(jsonentry2); err != nil {
		t.Error(err)
	}
	t.Log(e2)
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
