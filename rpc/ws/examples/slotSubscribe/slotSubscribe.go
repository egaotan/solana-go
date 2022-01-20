// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
	"os"
	"time"

	"github.com/gagliardetto/solana-go/rpc/ws"
)

func NewLog(dir, name string) *log.Logger {
	fileName := fmt.Sprintf("%s%s.log", dir, name)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log := log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	return log
}

func detect(name string, url string) {
	peers := make(map[string]bool, 0)
	for true {
		time.Sleep(time.Second * 30)
		//
		client, err := ws.Connect(context.Background(), url)
		if err != nil {
			continue
		}
		peer := client.Remote()
		_, ok := peers[peer]
		if !ok {
			peers[peer] = true
			go func(client *ws.Client) {
				defer client.Close()
				logger := NewLog("./", fmt.Sprintf("%s_%s", name, client.Remote()))
				sub, err := client.SlotSubscribe()
				if err != nil {
					panic(err)
				}
				defer sub.Unsubscribe()

				logger.Printf("node: %s", client.Remote())
				for {
					got, err := sub.Recv()
					if err != nil {
						panic(err)
					}
					logger.Printf("slot: %d", got.Slot)
				}
			}(client)
		}
	}
}

func main() {
	go detect("mainnet_beta", rpc.MainNetBeta_WS)
	go detect("serum", rpc.MainNetBetaSerum_WS)
	go detect("genesysgo", "wss://ssc-dao.genesysgo.net/")
	go detect("triton", "wss://autumn-empty-dawn.solana-mainnet.quiknode.pro/924b1527134b73309d1fd8b934a2f078ce31b189/")
	go detect("triton_free", "wss://free.rpcpool.com")
	time.Sleep(time.Hour * 2)
}
