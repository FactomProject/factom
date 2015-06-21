// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Balance struct {
	Balance int64
}

func ECBalance(key string) (*Balance, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/entry-credit-balance/%s", server, key))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
	b := new(Balance)
	if err := json.Unmarshal(body, b); err != nil {
		return nil, err
	}
	
	return b, nil
}

func FctBalance(key string) (*Balance, error) {
    str := fmt.Sprintf("http://%s/v1/factoid-balance/%s", serverFct, key)
    resp, err := http.Get(str)
    if err != nil {
        return nil, err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    resp.Body.Close()
        
    b := new(Balance)
    if err := json.Unmarshal(body, b); err != nil {
        return nil, err
    }
    
    return b, nil
}

func GenerateFactoidAddress(name string) (string, error) {
    
    type address struct {
        Address string
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-generate-address/%s", serverFct, name)
    resp, err := http.Get(str)
    if err != nil {
        return "", err
    }
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    resp.Body.Close()
    
    b := new(address)
    if err := json.Unmarshal(body, b); err != nil || len(b.Address)==0  {
        return "", fmt.Errorf("Duplicate or Invalid Name  ")
    }
    
    return b.Address, nil
}


func GenerateEntryCreditAddress(name string) (string, error) {
    fmt.Println("Here!!!!!!!!!!!!!!!!!!!!!!!")//////////////////////////////////////////////////////////////////
    type address struct {
        Address string
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-generate-ec-address/%s", serverFct, name)
    resp, err := http.Get(str)
    if err != nil {
        return "", err
    }
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    resp.Body.Close()
    
    b := new(address)
    if err := json.Unmarshal(body, b); err != nil || len(b.Address)==0  {
        return "", fmt.Errorf("Duplicate or Invalid Name  ")
    }
    
    return b.Address, nil
}