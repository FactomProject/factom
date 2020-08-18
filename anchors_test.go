package factom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAnchors(t *testing.T) {
	factomdResponse := `{"jsonrpc":"2.0","id":0,"result":{"directoryblockheight":200000,"directoryblockkeymr":"ce86fc790dd1462aea255adaa64e2f21c871995df2c2c119352d869fa1d7269f","bitcoin":{"transactionhash":"6d2d1e506528ae3b476d70fb05517bbbb152a4698a23ff78b4d87249027f53ca","blockhash":"0000000000000000000234e270b3fa6de63caad8a319731db5643ddb16b80cdf"},"ethereum":{"recordheight":200001,"dbheightmax":200000,"dbheightmin":199001,"windowmr":"935480547a2545161438da05136ec4238d88e7f8e1075687d8292fcafc1d0b22","merklebranch":[{"left":"a92a4460a9555c8b57a282e7ad4514d0f5aa8a612963da79db0e7cede6299bcd","right":"ce86fc790dd1462aea255adaa64e2f21c871995df2c2c119352d869fa1d7269f","top":"bbbf494fff20d8fbada47980498a96e49fba920fd1ca4e0586af4961e4318354"},{"left":"b41ffe5f9b710edf6bcf1a6d3e92685155411111d640f2b2224ab15f04c46161","right":"bbbf494fff20d8fbada47980498a96e49fba920fd1ca4e0586af4961e4318354","top":"ca1e7d59dc295c48ebddf6c068826584fb5adb38798f5cac7f964509627825f5"},{"left":"c919c0fb94900f3016c75a7f9f992c79aa7db126db6c6c3b1af92bbee174bda6","right":"ca1e7d59dc295c48ebddf6c068826584fb5adb38798f5cac7f964509627825f5","top":"3c8e97e32b4dbc0e4d82f066217492ca72ca722b91dbc764c311c39a9196bdda"},{"left":"3c8e97e32b4dbc0e4d82f066217492ca72ca722b91dbc764c311c39a9196bdda","right":"3c8e97e32b4dbc0e4d82f066217492ca72ca722b91dbc764c311c39a9196bdda","top":"5713abb24c8a8480f5c04bfdb8668a93e8b0a24a0562afbc7f362d665938a89e"},{"left":"5713abb24c8a8480f5c04bfdb8668a93e8b0a24a0562afbc7f362d665938a89e","right":"5713abb24c8a8480f5c04bfdb8668a93e8b0a24a0562afbc7f362d665938a89e","top":"a7ab3a72f5e754f3e10f04b9f7f2e50f36c688776b521918a755e62fd6c2da34"},{"left":"5e42d2f5ba59338b317369bbc53fbb70c094364ae8c15c9a291dad65dec9d839","right":"a7ab3a72f5e754f3e10f04b9f7f2e50f36c688776b521918a755e62fd6c2da34","top":"bc5c336d7e53b6d9f2832038c69a78eb7300bd539b6002d44f450a4d2c979aa3"},{"left":"31041410eb7bd6a9fabe974a0f4eabe477ffbbf5806fdaee5da98faff63b8142","right":"bc5c336d7e53b6d9f2832038c69a78eb7300bd539b6002d44f450a4d2c979aa3","top":"f9b5967e2acadfac93bc2141ee9f85b0649f89452e44e4f8b1cbe877457a0da4"},{"left":"04aa1b2d35ea99becd92bcfccc7406aea37287a58425b83349d0fc8cc7bec443","right":"f9b5967e2acadfac93bc2141ee9f85b0649f89452e44e4f8b1cbe877457a0da4","top":"77eca087628e243bf2914162be6c4aefaf736ed56e8c125b955cef3c48a1ab4b"},{"left":"be8394e3f3a205ad37afb8af778695e2914f071b3d215e586a93288fd692eb00","right":"77eca087628e243bf2914162be6c4aefaf736ed56e8c125b955cef3c48a1ab4b","top":"12aa5f990bd209dfaa83d9ea050e96e3b44e54ffee02252028fe969da3ab690d"},{"left":"2bdead51afe10d07c5dc5b29351fedd4a13849261359a6f22100dcc7ff50543a","right":"12aa5f990bd209dfaa83d9ea050e96e3b44e54ffee02252028fe969da3ab690d","top":"935480547a2545161438da05136ec4238d88e7f8e1075687d8292fcafc1d0b22"}],"contractaddress":"0xfac701d9554a008e48b6307fb90457ba3959e8a8","txid":"0xd57492c9d505e3052c454acdfc3768bc3eb8859c91829654346503ef2dcb6a23","blockhash":"0x49e8f3c394079e6c3964fdf943ea9759475d4ca5aec90c9a695be07efdd88d32","txindex":31}}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, err := GetAnchors("ce86fc790dd1462aea255adaa64e2f21c871995df2c2c119352d869fa1d7269f") // irrelevant, hardcoded above
	if err != nil {
		t.Error(err)
	}
	response2, err := GetAnchorsByHeight(200000) // irrelevant, hardcoded above
	if err != nil {
		t.Error(err)
	}

	received1 := fmt.Sprintf("%+v", response)
	received2 := fmt.Sprintf("%+v", response2)
	expected := `Height: 200000
KeyMR: ce86fc790dd1462aea255adaa64e2f21c871995df2c2c119352d869fa1d7269f
Bitcoin {
 TransactionHash: 6d2d1e506528ae3b476d70fb05517bbbb152a4698a23ff78b4d87249027f53ca
 BlockHash: 0000000000000000000234e270b3fa6de63caad8a319731db5643ddb16b80cdf
}
Ethereum {
 RecordHeight: 200001
 DBHeightMax: 200000
 DBHeightMin: 199001
 WindowMR: 935480547a2545161438da05136ec4238d88e7f8e1075687d8292fcafc1d0b22
  MerkleBranch {
   Branch {
    Left: a92a4460a9555c8b57a282e7ad4514d0f5aa8a612963da79db0e7cede6299bcd
    Right: ce86fc790dd1462aea255adaa64e2f21c871995df2c2c119352d869fa1d7269f
    Top: bbbf494fff20d8fbada47980498a96e49fba920fd1ca4e0586af4961e4318354
   Branch }
   Branch {
    Left: b41ffe5f9b710edf6bcf1a6d3e92685155411111d640f2b2224ab15f04c46161
    Right: bbbf494fff20d8fbada47980498a96e49fba920fd1ca4e0586af4961e4318354
    Top: ca1e7d59dc295c48ebddf6c068826584fb5adb38798f5cac7f964509627825f5
   Branch }
   Branch {
    Left: c919c0fb94900f3016c75a7f9f992c79aa7db126db6c6c3b1af92bbee174bda6
    Right: ca1e7d59dc295c48ebddf6c068826584fb5adb38798f5cac7f964509627825f5
    Top: 3c8e97e32b4dbc0e4d82f066217492ca72ca722b91dbc764c311c39a9196bdda
   Branch }
   Branch {
    Left: 3c8e97e32b4dbc0e4d82f066217492ca72ca722b91dbc764c311c39a9196bdda
    Right: 3c8e97e32b4dbc0e4d82f066217492ca72ca722b91dbc764c311c39a9196bdda
    Top: 5713abb24c8a8480f5c04bfdb8668a93e8b0a24a0562afbc7f362d665938a89e
   Branch }
   Branch {
    Left: 5713abb24c8a8480f5c04bfdb8668a93e8b0a24a0562afbc7f362d665938a89e
    Right: 5713abb24c8a8480f5c04bfdb8668a93e8b0a24a0562afbc7f362d665938a89e
    Top: a7ab3a72f5e754f3e10f04b9f7f2e50f36c688776b521918a755e62fd6c2da34
   Branch }
   Branch {
    Left: 5e42d2f5ba59338b317369bbc53fbb70c094364ae8c15c9a291dad65dec9d839
    Right: a7ab3a72f5e754f3e10f04b9f7f2e50f36c688776b521918a755e62fd6c2da34
    Top: bc5c336d7e53b6d9f2832038c69a78eb7300bd539b6002d44f450a4d2c979aa3
   Branch }
   Branch {
    Left: 31041410eb7bd6a9fabe974a0f4eabe477ffbbf5806fdaee5da98faff63b8142
    Right: bc5c336d7e53b6d9f2832038c69a78eb7300bd539b6002d44f450a4d2c979aa3
    Top: f9b5967e2acadfac93bc2141ee9f85b0649f89452e44e4f8b1cbe877457a0da4
   Branch }
   Branch {
    Left: 04aa1b2d35ea99becd92bcfccc7406aea37287a58425b83349d0fc8cc7bec443
    Right: f9b5967e2acadfac93bc2141ee9f85b0649f89452e44e4f8b1cbe877457a0da4
    Top: 77eca087628e243bf2914162be6c4aefaf736ed56e8c125b955cef3c48a1ab4b
   Branch }
   Branch {
    Left: be8394e3f3a205ad37afb8af778695e2914f071b3d215e586a93288fd692eb00
    Right: 77eca087628e243bf2914162be6c4aefaf736ed56e8c125b955cef3c48a1ab4b
    Top: 12aa5f990bd209dfaa83d9ea050e96e3b44e54ffee02252028fe969da3ab690d
   Branch }
   Branch {
    Left: 2bdead51afe10d07c5dc5b29351fedd4a13849261359a6f22100dcc7ff50543a
    Right: 12aa5f990bd209dfaa83d9ea050e96e3b44e54ffee02252028fe969da3ab690d
    Top: 935480547a2545161438da05136ec4238d88e7f8e1075687d8292fcafc1d0b22
   Branch }
  }
 ContractAddress: 0xfac701d9554a008e48b6307fb90457ba3959e8a8
 TxID: 0xd57492c9d505e3052c454acdfc3768bc3eb8859c91829654346503ef2dcb6a23
 BlockHash: 0x49e8f3c394079e6c3964fdf943ea9759475d4ca5aec90c9a695be07efdd88d32
 TxIndex: 31
}
`
	if received1 != expected {
		t.Errorf("GetAnchors() expected:%s\nreceived:%s", expected, received1)
	}
	if received2 != expected {
		t.Errorf("GetAnchorsByHeight() expected:%s\nreceived:%s", expected, received2)
	}

}
