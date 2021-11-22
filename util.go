package main

import (
	"encoding/json"
	"fmt"

	"github.com/harmony-one/go-sdk/pkg/rpc"

	"github.com/harmony-one/go-sdk/pkg/common"
)

var request = func(method string, params []interface{}) error {
	success, failure := rpc.Request(method, node, params)
	if failure != nil {
		return failure
	}
	asJSON, _ := json.Marshal(success)

	fmt.Println(common.JSONPrettyFormat(string(asJSON)))
	return nil
}
