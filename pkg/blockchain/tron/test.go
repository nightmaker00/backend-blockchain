// package main

// import (
// 	"fmt"
// 	"strings"
// 	"net/http"
// 	"io"
// )

// func main() {

// 	url := "https://api.shasta.trongrid.io/wallet/createtransaction"

// 	payload := strings.NewReader("{\"owner_address\":\"TZ4UXDV5ZhNW7fb2AMSbgfAEZ7hWsnYS2g\",\"to_address\":\"TPswDDCAWhJAZGdHPidFg5nEf8TkNToDX1\",\"amount\":1000,\"visible\":true}")

// 	req, _ := http.NewRequest("POST", url, payload)

// 	req.Header.Add("accept", "application/json")
// 	req.Header.Add("content-type", "application/json")

// 	res, _ := http.DefaultClient.Do(req)

// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)

// 	fmt.Println(string(body))

// }

package tron