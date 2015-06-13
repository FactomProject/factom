package factom

import (
	"fmt"
	"net/http"
)

func TestCredit(eckey string) error {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/test-credit/%s", server, eckey))
	if err != nil {
		return err
	}
	resp.Body.Close()
	
	return nil
}