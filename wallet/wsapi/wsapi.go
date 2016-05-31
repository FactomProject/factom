// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

import (
	"encoding/json"
	"io/ioutil"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/web"
)

const API_VERSION string = "2.0"

var (
	webServer *web.Server
	fctWallet *wallet.Wallet
)

func Start(w *wallet.Wallet, net string) {
	webServer = web.NewServer()
	fctWallet = w

	webServer.Post("/v2", handleV2)
	webServer.Get("/v2", handleV2)
	webServer.Run(net)
}

func Stop() {
	fctWallet.Close()
	webServer.Close()
}

func handleV2(ctx *web.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleV2Error(ctx, nil, newInvalidRequestError())
		return
	}

	j, err := factom.ParseJSON2Request(string(body))
	if err != nil {
		handleV2Error(ctx, nil, newInvalidRequestError())
		return
	}

	jsonResp, jsonError := handleV2Request(j)

	if jsonError != nil {
		handleV2Error(ctx, j, jsonError)
		return
	}

	ctx.Write([]byte(jsonResp.String()))
}

func handleV2Request(j *factom.JSON2Request) (*factom.JSON2Response, *factom.JSONError) {
	var resp interface{}
	var jsonError *factom.JSONError
	params := j.Params
	switch j.Method {
	case "test":
		resp, jsonError = handleTest(params)
	case "address":
		resp, jsonError = handleAddress(params)
	case "all-addresses":
		resp, jsonError = handleAllAddresses(params)
	case "generate-ec-address":
		resp, jsonError = handleGenerateECAddress(params)
	case "generate-factoid-address":
		resp, jsonError = handleGenerateFactoidAddress(params)
	default:
		jsonError = newMethodNotFoundError()
	}
	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := factom.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}

func handleTest(params interface{}) (interface{}, *factom.JSONError) {
	return "Hello Factom!", nil
}

func handleAddress(params interface{}) (interface{}, *factom.JSONError) {
	req := new(addressRequest)
	if err := mapToObject(params, req); err != nil {
		return nil, newInvalidParamsError()
	}

	resp := new(addressResponse)
	if e, err := fctWallet.GetECAddress(req.Address); err != nil {
		return nil, newCustomInternalError(err)
	} else if e != nil {
		resp.Public = e.PubString()
		resp.Secret = e.SecString()
	}
	if f, err := fctWallet.GetFCTAddress(req.Address); err != nil {
		return nil, newCustomInternalError(err)
	} else if f != nil {
		resp.Public = f.PubString()
		resp.Secret = f.SecString()
	}
	if resp.Secret == "" {
		return nil, newCustomInternalError("No Addresses Found")
	}

	return resp, nil
}

func handleAllAddresses(params interface{}) (interface{}, *factom.JSONError) {
	resp := new(allAddressesResponse)

	fs, es, err := fctWallet.GetAllAddresses()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	for _, f := range fs {
		a := new(addressResponse)
		a.Public = f.PubString()
		a.Secret = f.SecString()
		resp.Addresses = append(resp.Addresses, a)
	}
	for _, e := range es {
		a := new(addressResponse)
		a.Public = e.PubString()
		a.Secret = e.SecString()
		resp.Addresses = append(resp.Addresses, a)
	}

	return resp, nil
}

func handleGenerateFactoidAddress(params interface{}) (interface{}, *factom.JSONError) {
	a, err := fctWallet.GenerateFCTAddress()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	
	resp := new(addressResponse)
	resp.Public = a.PubString()
	resp.Secret = a.SecString()
	
	return resp, nil
}

func handleGenerateECAddress(params interface{}) (interface{}, *factom.JSONError) {
	a, err := fctWallet.GenerateECAddress()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	
	resp := new(addressResponse)
	resp.Public = a.PubString()
	resp.Secret = a.SecString()
	
	return resp, nil
}

// utility functions

func mapToObject(source interface{}, dst interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
