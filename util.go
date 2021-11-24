package main

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
)

var request = func(method string, params []interface{}) (result rpc.Reply, err error) {
	result, err = rpc.Request(method, node, params)
	return
}
