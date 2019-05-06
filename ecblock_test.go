// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"encoding/json"
	"fmt"

	. "github.com/FactomProject/factom"
)

func TestUnmarshalECBlock(t *testing.T) {
	js := []byte(`{"ecblock":{"header":{"bodyhash":"541338744c8254641e0df2776dc7af07915c5da009e72e764da2bcbaa29a1bc6","prevheaderhash":"86aa9a8ef0cdb5e7b525fb7f9dd05f8188471cfbea6cf1c7ebab482ec408b6e9","prevfullhash":"af8a96d6e4ce0bd81c327bc49ab96c7e190c08c5ea0257d95a88c0806abf4266","dbheight":10199,"headerexpansionarea":"","objectcount":14,"bodysize":561,"chainid":"000000000000000000000000000000000000000000000000000000000000000c","ecchainid":"000000000000000000000000000000000000000000000000000000000000000c"},"body":{"entries":[{"serverindexnumber":0},{"version":0,"millitime":"0150f7d966a9","chainidhash":"e5f6f7cd369ef90a9872532af2d9755edfcd78124ea140f3417f54949b169aea","weld":"1aa415bfaa978342ef396d7203cde3ad45cf92dab89ec6b34128234cae42ef6f","entryhash":"7b4bc033547fd3ac1055d500752e99048d83ae9e580cc1fa4dcead10db868c73","credits":11,"ecpubkey":"79a1ad273d890287e5d4f16d2669c06c523b9e48673de1bfde3ea2fda309ac92","sig":"34cab18fbc270bc51e9d68adc8cb9c65da5d7021bcc34370598ac6370fb7edde9b5c1a0164055bef53a83fbb1ddeb61a6942491fd8f9a56eb264c1abcc7c3905"},{"version":0,"millitime":"0150f7d8f870","entryhash":"ac43f66ddf733981ce33a15bff872e125fff1a2b640cf99ee7e44b6ca2e96fb6","credits":1,"ecpubkey":"4bcbc1c5ab90e432bd407a51eaa513b4050eecda1fd42bbf6b7050a1d96f94b7","sig":"d06dedddf728f55a011eb6c133bfeebe1669823afd109158f9c6cbeaf012d358e9bc0055850ca639bb78838418465e48aa1f9e03874c948e8520d9064adb9c06"},{"number":1},{"number":2},{"number":3},{"number":4},{"version":0,"millitime":"0150f7dcfb53","chainidhash":"1962219a271a272ff432fb8635ce07269d6f4a974871bbfde9d5ac7ab429a682","weld":"2b5088c89e158f94802459c01a9eb170eca3487f4de26ff8a331a5b5f5dbde4e","entryhash":"8c138dfb419a2c118c58a7ac0e791c3c6c2a67cec732325c2465ce911af41a4e","credits":11,"ecpubkey":"79a1ad273d890287e5d4f16d2669c06c523b9e48673de1bfde3ea2fda309ac92","sig":"ff2a6878ab59da88bd15b94545fbdecbab29fd14f64e7d7cf5fe3eb7f2f08a169aa1cfea415bd5d86d934ff925dfd8567491bdc7d9dff2a38d28bed729364101"},{"number":5},{"number":6},{"number":7},{"number":8},{"number":9},{"number":10}]},"headerhash":"a7baaa24e477a0acef165461d70ec94ff3f33ad15562ecbe937967a761929a17","fullhash":"84339a4a849c3616c7c1a5011f2fe14d000efd3a98309afaabbd2d7c0122094c"},"rawdata":"000000000000000000000000000000000000000000000000000000000000000c541338744c8254641e0df2776dc7af07915c5da009e72e764da2bcbaa29a1bc686aa9a8ef0cdb5e7b525fb7f9dd05f8188471cfbea6cf1c7ebab482ec408b6e9af8a96d6e4ce0bd81c327bc49ab96c7e190c08c5ea0257d95a88c0806abf4266000027d700000000000000000e0000000000000231000002000150f7d966a9e5f6f7cd369ef90a9872532af2d9755edfcd78124ea140f3417f54949b169aea1aa415bfaa978342ef396d7203cde3ad45cf92dab89ec6b34128234cae42ef6f7b4bc033547fd3ac1055d500752e99048d83ae9e580cc1fa4dcead10db868c730b79a1ad273d890287e5d4f16d2669c06c523b9e48673de1bfde3ea2fda309ac9234cab18fbc270bc51e9d68adc8cb9c65da5d7021bcc34370598ac6370fb7edde9b5c1a0164055bef53a83fbb1ddeb61a6942491fd8f9a56eb264c1abcc7c390503000150f7d8f870ac43f66ddf733981ce33a15bff872e125fff1a2b640cf99ee7e44b6ca2e96fb6014bcbc1c5ab90e432bd407a51eaa513b4050eecda1fd42bbf6b7050a1d96f94b7d06dedddf728f55a011eb6c133bfeebe1669823afd109158f9c6cbeaf012d358e9bc0055850ca639bb78838418465e48aa1f9e03874c948e8520d9064adb9c06010101020103010402000150f7dcfb531962219a271a272ff432fb8635ce07269d6f4a974871bbfde9d5ac7ab429a6822b5088c89e158f94802459c01a9eb170eca3487f4de26ff8a331a5b5f5dbde4e8c138dfb419a2c118c58a7ac0e791c3c6c2a67cec732325c2465ce911af41a4e0b79a1ad273d890287e5d4f16d2669c06c523b9e48673de1bfde3ea2fda309ac92ff2a6878ab59da88bd15b94545fbdecbab29fd14f64e7d7cf5fe3eb7f2f08a169aa1cfea415bd5d86d934ff925dfd8567491bdc7d9dff2a38d28bed72936410101050106010701080109010a"}`)

	jsbadentry := []byte(`{"ecblock":{"header":{"bodyhash":"541338744c8254641e0df2776dc7af07915c5da009e72e764da2bcbaa29a1bc6","prevheaderhash":"86aa9a8ef0cdb5e7b525fb7f9dd05f8188471cfbea6cf1c7ebab482ec408b6e9","prevfullhash":"af8a96d6e4ce0bd81c327bc49ab96c7e190c08c5ea0257d95a88c0806abf4266","dbheight":10199,"headerexpansionarea":"","objectcount":14,"bodysize":561,"chainid":"000000000000000000000000000000000000000000000000000000000000000c","ecchainid":"000000000000000000000000000000000000000000000000000000000000000c"},"body":{"entries":[{"badentry":"bad"},{"serverindexnumber":0},{"number":5},{"number":6},{"number":7},{"number":8},{"number":9},{"number":10}]},"headerhash":"a7baaa24e477a0acef165461d70ec94ff3f33ad15562ecbe937967a761929a17","fullhash":"84339a4a849c3616c7c1a5011f2fe14d000efd3a98309afaabbd2d7c0122094c"},"rawdata":"000000000000000000000000000000000000000000000000000000000000000c541338744c8254641e0df2776dc7af07915c5da009e72e764da2bcbaa29a1bc686aa9a8ef0cdb5e7b525fb7f9dd05f8188471cfbea6cf1c7ebab482ec408b6e9af8a96d6e4ce0bd81c327bc49ab96c7e190c08c5ea0257d95a88c0806abf4266000027d700000000000000000e0000000000000231000002000150f7d966a9e5f6f7cd369ef90a9872532af2d9755edfcd78124ea140f3417f54949b169aea1aa415bfaa978342ef396d7203cde3ad45cf92dab89ec6b34128234cae42ef6f7b4bc033547fd3ac1055d500752e99048d83ae9e580cc1fa4dcead10db868c730b79a1ad273d890287e5d4f16d2669c06c523b9e48673de1bfde3ea2fda309ac9234cab18fbc270bc51e9d68adc8cb9c65da5d7021bcc34370598ac6370fb7edde9b5c1a0164055bef53a83fbb1ddeb61a6942491fd8f9a56eb264c1abcc7c390503000150f7d8f870ac43f66ddf733981ce33a15bff872e125fff1a2b640cf99ee7e44b6ca2e96fb6014bcbc1c5ab90e432bd407a51eaa513b4050eecda1fd42bbf6b7050a1d96f94b7d06dedddf728f55a011eb6c133bfeebe1669823afd109158f9c6cbeaf012d358e9bc0055850ca639bb78838418465e48aa1f9e03874c948e8520d9064adb9c06010101020103010402000150f7dcfb531962219a271a272ff432fb8635ce07269d6f4a974871bbfde9d5ac7ab429a6822b5088c89e158f94802459c01a9eb170eca3487f4de26ff8a331a5b5f5dbde4e8c138dfb419a2c118c58a7ac0e791c3c6c2a67cec732325c2465ce911af41a4e0b79a1ad273d890287e5d4f16d2669c06c523b9e48673de1bfde3ea2fda309ac92ff2a6878ab59da88bd15b94545fbdecbab29fd14f64e7d7cf5fe3eb7f2f08a169aa1cfea415bd5d86d934ff925dfd8567491bdc7d9dff2a38d28bed72936410101050106010701080109010a"}`)
	wrap := new(struct {
		ECBlock ECBlock `json:"ecblock"`
		RawData string  `json:"rawdata"`
	})

	err := json.Unmarshal(js, wrap)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(jsbadentry, wrap)
	if err != ErrUnknownECBEntry {
		t.Error(err)
	}

	t.Log("ECBlock:", wrap.ECBlock)
	t.Log("RawData:", wrap.RawData)
}

func TestGetECBlock(t *testing.T) {
	// Check for a missing blockHash
	_, _, err := GetECBlock("baadbaadbaadbaadbaadbaadbaadbaadbaadbaadbaadbaadbaadbaadbaadbaad")
	if err == nil {
		t.Error("expected error for missing block")
	} else {
		t.Log("Missing Block Error:", err)
	}

	ecb, raw, err := GetECBlock("a7baaa24e477a0acef165461d70ec94ff3f33ad15562ecbe937967a761929a17")
	if err != nil {
		t.Error(err)
	}
	t.Log("ECBlock: ", ecb)
	t.Log(fmt.Sprintf("raw: %x\n", raw))
}

func TestGetECBlockByHeight(t *testing.T) {
	ecb, raw, err := GetECBlockByHeight(10199)
	if err != nil {
		t.Error(err)
	}
	t.Log("ECBlock: ", ecb)
	t.Log(fmt.Sprintf("raw: %x\n", raw))
}
