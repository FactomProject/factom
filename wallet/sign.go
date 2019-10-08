package wallet

import (
	"github.com/FactomProject/factomd/common/primitives"
)

// SignData signs arbitrary data
func (w *Wallet) SignData(signer string, data []byte) ([]byte, []byte, error) {

	var priv []byte
	var pub []byte

	if fa, err := w.GetFCTAddress(signer); err == nil {
		priv = fa.SecBytes()
		pub = fa.PubBytes()
	} else if ec, err := w.GetECAddress(signer); err == nil {
		priv = ec.SecBytes()
		pub = ec.PubBytes()
	} else if id, err := w.GetIdentityKey(signer); err == nil {
		priv = id.SecBytes()
		pub = id.PubBytes()
	} else {
		return nil, nil, ErrNoSuchAddress
	}

	sig := primitives.Sign(priv, data)
	return pub, sig, nil
}
