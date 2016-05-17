// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

import (
	"io/ioutil"
	
	"github.com/FactomProject/factom"
	"github.com/FactomProject/web"
)

const API_VERSION string = "2.0"

var server = web.NewServer()

func Start(net string) {
	server.Post("/v2", HandleV2)
	server.Get("/v2", HandleV2)
	server.Run(net)
}

func Stop() {
	server.Close()
}

func HandleV2(ctx *web.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		HandleV2Error(ctx, nil, NewInvalidRequestError())
		return
	}

	j, err := factom.ParseJSON2Request(string(body))
	if err != nil {
		HandleV2Error(ctx, nil, NewInvalidRequestError())
		return
	}

	jsonResp, jsonError := HandleV2Request(j)

	if jsonError != nil {
		HandleV2Error(ctx, j, jsonError)
		return
	}

	ctx.Write([]byte(jsonResp.String()))
}

func HandleV2Request(j *factom.JSON2Request) (*factom.JSON2Response, *factom.JSONError) {
	var resp interface{}
	var jsonError *factom.JSONError
	params := j.Params
	switch j.Method {
	case "test":
		resp, jsonError = HandleTest(params)
	default:
		jsonError = NewMethodNotFoundError()
	}
	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := factom.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}

func HandleTest(params interface{}) (interface{}, *factom.JSONError) {
	return "Hello Factom!", nil
}
