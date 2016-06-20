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
	
	"fmt" // DEBUG
)

var _ = fmt.Sprint("DEBUG")

const APIVersion string = "2.0"

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
	params := []byte(j.Params)

	switch j.Method {
	case "address":
		resp, jsonError = handleAddress(params)
	case "all-addresses":
		resp, jsonError = handleAllAddresses(params)
	case "generate-ec-address":
		resp, jsonError = handleGenerateECAddress(params)
	case "generate-factoid-address":
		resp, jsonError = handleGenerateFactoidAddress(params)
	case "import-addresses":
		resp, jsonError = handleImportAddresses(params)
	case "wallet-backup":
		resp, jsonError = handleWalletBackup(params)
	case "new-transaction":
		resp, jsonError = handleNewTransaction(params)
	case "delete-transaction":
		resp, jsonError = handleDeleteTransaction(params)
	case "transactions":
		resp, jsonError = handleTransactions(params)
	case "add-input":
		resp, jsonError = handleAddInput(params)
	case "add-output":
		resp, jsonError = handleAddOutput(params)
	case "add-ec-output":
		resp, jsonError = handleAddECOutput(params)
	case "add-fee":
		resp, jsonError = handleAddFee(params)
	case "sub-fee":
		resp, jsonError = handleSubFee(params)
	case "sign-transaction":
		resp, jsonError = handleSignTransaction(params)
	case "compose-transaction":
		resp, jsonError = handleComposeTransaction(params)
	default:
		jsonError = newMethodNotFoundError()
	}
	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := factom.NewJSON2Response()
	jsonResp.ID = j.ID
	if b, err := json.Marshal(resp); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		jsonResp.Result = b
	}

	return jsonResp, nil
}

func handleAddress(params []byte) (interface{}, *factom.JSONError) {
	req := new(addressRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}

	resp := new(addressResponse)
	switch factom.AddressStringType(req.Address) {
	case factom.ECPub:
		e, err := fctWallet.GetECAddress(req.Address)
		if err != nil {
			return nil, newCustomInternalError(err)
		}
		resp = mkAddressResponse(e)
	case factom.FactoidPub:
		f, err := fctWallet.GetFCTAddress(req.Address)
		if err != nil {
			return nil, newCustomInternalError(err)
		}
		resp = mkAddressResponse(f)
	default:
		return nil, newCustomInternalError("Invalid address type")
	}

	return resp, nil
}

func handleAllAddresses(params []byte) (interface{}, *factom.JSONError) {
	resp := new(multiAddressResponse)

	fs, es, err := fctWallet.GetAllAddresses()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	for _, f := range fs {
		a := mkAddressResponse(f)
		resp.Addresses = append(resp.Addresses, a)
	}
	for _, e := range es {
		a := mkAddressResponse(e)
		resp.Addresses = append(resp.Addresses, a)
	}

	return resp, nil
}

func handleGenerateFactoidAddress(params []byte) (interface{}, *factom.JSONError) {
	a, err := fctWallet.GenerateFCTAddress()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	
	resp := mkAddressResponse(a)
	
	return resp, nil
}

func handleGenerateECAddress(params []byte) (interface{}, *factom.JSONError) {
	a, err := fctWallet.GenerateECAddress()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	
	resp := mkAddressResponse(a)
	
	return resp, nil
}

func handleImportAddresses(params []byte)  (interface{}, *factom.JSONError) {
	req := new(importRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	resp := new(multiAddressResponse)
	for _, v := range req.Addresses {
		switch factom.AddressStringType(v.Secret) {
		case factom.FactoidSec:
			f, err := factom.GetFactoidAddress(v.Secret)
			if err != nil {
				return nil, newCustomInternalError(err)
			}
			if err := fctWallet.PutFCTAddress(f); err != nil {
				return nil, newCustomInternalError(err)
			}
			a := mkAddressResponse(f)
			resp.Addresses = append(resp.Addresses, a)
		case factom.ECSec:
			e, err := factom.GetECAddress(v.Secret)
			if err != nil {
				return nil, newCustomInternalError(err)
			}
			if err := fctWallet.PutECAddress(e); err != nil {
				return nil, newCustomInternalError(err)
			}
			a := mkAddressResponse(e)
			resp.Addresses = append(resp.Addresses, a)
		default:
			return nil, newCustomInternalError("address could not be imported")
		}
	}
	return resp, nil
}

func handleWalletBackup(params []byte) (interface{}, *factom.JSONError) {
	resp := new(walletBackupResponse)

	if seed, err := fctWallet.GetSeed(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Seed = seed
	}
	
	fs, es, err := fctWallet.GetAllAddresses()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	for _, f := range fs {
		a := mkAddressResponse(f)
		resp.Addresses = append(resp.Addresses, a)
	}
	for _, e := range es {
		a := mkAddressResponse(e)
		resp.Addresses = append(resp.Addresses, a)
	}

	return resp, nil
}

// transaction handlers

func handleNewTransaction(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	if err := fctWallet.NewTransaction(req.Name); err != nil {
		return nil, newCustomInternalError(err)
	}

	resp := transactionResponse{Name: req.Name}
	t := fctWallet.GetTransactions()[req.Name]
	if s, err := t.JSONString(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Transaction = s
	}
	
	return resp, nil
}

func handleDeleteTransaction(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	if err := fctWallet.DeleteTransaction(req.Name); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	return resp, nil
}

func handleTransactions(params []byte) (interface{}, *factom.JSONError) {
	resp := new(transactionsResponse)

	for name, _ := range fctWallet.GetTransactions() {
		r := transactionResponse{Name: name}
		t := fctWallet.GetTransactions()[name]
		if s, err := t.JSONString(); err != nil {
			return nil, newCustomInternalError(err)
		} else {
			r.Transaction = s
		}
		resp.Transactions = append(resp.Transactions, r)
	}
	
	return resp, nil
}

func handleAddInput(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionValueRequest)
	if err := json.Unmarshal(params, req); err != nil {
		fmt.Println("DEBUG:", err)
		return nil, newInvalidParamsError()
	}
	
	if err := fctWallet.AddInput(req.Name, req.Address, req.Amount); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	t := fctWallet.GetTransactions()[req.Name]
	if s, err := t.JSONString(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Transaction = s
	}
	
	return resp, nil
}

func handleAddOutput(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionValueRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	if err := fctWallet.AddOutput(req.Name, req.Address, req.Amount); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	return resp, nil
}

func handleAddECOutput(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionValueRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	if err := fctWallet.AddECOutput(req.Name, req.Address, req.Amount); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	t := fctWallet.GetTransactions()[req.Name]
	if s, err := t.JSONString(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Transaction = s
	}
	
	return resp, nil
}

func handleAddFee(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionAddressRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	rate, err := factom.GetRate()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	if err := fctWallet.AddFee(req.Name, req.Address, rate); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	t := fctWallet.GetTransactions()[req.Name]
	if s, err := t.JSONString(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Transaction = s
	}
	
	return resp, nil
}

func handleSubFee(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionAddressRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	rate, err := factom.GetRate()
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	if err := fctWallet.SubFee(req.Name, req.Address, rate); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	t := fctWallet.GetTransactions()[req.Name]
	if s, err := t.JSONString(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Transaction = s
	}
	
	return resp, nil
}

func handleSignTransaction(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	if err := fctWallet.SignTransaction(req.Name); err != nil {
		return nil, newCustomInternalError(err)
	}
	resp := transactionResponse{Name: req.Name}
	t := fctWallet.GetTransactions()[req.Name]
	if s, err := t.JSONString(); err != nil {
		return nil, newCustomInternalError(err)
	} else {
		resp.Transaction = s
	}
	
	return resp, nil
}

func handleComposeTransaction(params []byte) (interface{}, *factom.JSONError) {
	req := new(transactionRequest)
	if err := json.Unmarshal(params, req); err != nil {
		return nil, newInvalidParamsError()
	}
	
	t, err := fctWallet.ComposeTransaction(req.Name)
	if err != nil {
		return nil, newCustomInternalError(err)
	}
	return t, nil
}

// utility functions

type addressResponder interface {
	String() string
	SecString() string
}

func mkAddressResponse(a addressResponder) *addressResponse {
	r := new(addressResponse)
	r.Public = a.String()
	r.Secret = a.SecString()
	return r
}
