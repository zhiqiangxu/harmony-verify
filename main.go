package main

import "github.com/harmony-one/go-sdk/pkg/rpc"

const node = "https://api.s1.b.hmny.io"

func main() {
	request(rpc.Method.GetLatestBlockHeader, []interface{}{})
}
