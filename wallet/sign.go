package wallet

import (
	"github.com/FactomProject/factomd/common/primitives"
)

// SignData signs arbitrary data
func (w *Wallet) SignData(address string, data []byte) ([]byte, []byte, error) {
	f, err := w.GetFCTAddress(address)
	if err != nil {
		return nil, nil, err
	}

	sig := primitives.Sign(f.SecBytes(), data)
	return f.PubBytes(), sig, nil
}
