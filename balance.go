// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
    "strings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)



func ECBalance(key string) (int64, error) {
    str := fmt.Sprintf("http://%s/v1/entry-credit-balance/%s", serverFct, key)
    resp, err := http.Get(str)
	if err != nil {
		return 0, err
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
    
    type Balance struct {
        Balance int64
    }	
    
    b := new(Balance)
	if err := json.Unmarshal(body, b); err != nil {
        return 0, err
	}
	
	return b.Balance, nil
}

func FctBalance(key string) (int64, error) {
    str := fmt.Sprintf("http://%s/v1/factoid-balance/%s", serverFct, key)
    resp, err := http.Get(str)
    if err != nil {
        return 0, err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, err
    }
    resp.Body.Close()
        
    type Balance struct {
        Balance int64
    }
    
    b := new(Balance)
    if err := json.Unmarshal(body, b); err != nil {
        return 0, err
    }
    
    return b.Balance, nil
}

func GenerateFactoidAddress(name string) (string, error) {
    name = strings.TrimSpace(name)
    
    type address struct {
        Address string
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-generate-address/%s", serverFct, name)
    fmt.Println(str)/////////////////////////////////////////////////////////////////////////////////////////
    
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
        return "", fmt.Errorf("Factoid Address Generation Failed: "+string(body))
    }
    
    return b.Address, nil
}


func GenerateEntryCreditAddress(name string) (string, error) {
    name = strings.TrimSpace(name)
    
    type address struct {
        Address string
        Success bool
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-generate-ec-address/%s", serverFct, name)
    fmt.Println(str)//////////////////////////////////////////////////////////////////////////////
    

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
        return "", fmt.Errorf("Entry Credit Address Generation Failed: "+string(body))
    }

    return b.Address, nil
}