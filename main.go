package main

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"

	"encoding/json"
	"fmt"
	"github.com/harmony-one/go-sdk/pkg/common"
)

const node = "https://api.s1.b.hmny.io"

func main() {
	{
		result, err := request(rpc.Method.GetLatestBlockHeader, []interface{}{})
		if err != nil {
			panic(err)
		}

		asJSON, _ := json.Marshal(result)
		fmt.Println(common.JSONPrettyFormat(string(asJSON)))
	}

	{
		result, err := request(rpc.Method.GetBlockByNumber, []interface{}{int64(10)})
		if err != nil {
			panic(err)
		}

		asJSON, _ := json.Marshal(result)
		fmt.Println(common.JSONPrettyFormat(string(asJSON)))
	}
}
