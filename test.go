package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	type RemoteRequest struct {
		Req string
	}

	data := "{\"Req\":\"POST / HTTP/1.1\\r\\nHost: localhost:8080\\r\\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36\\r\\nContent-Length: 94\\r\\nAccept: */*\\r\\nAccept-Encoding: gzip, deflate, br\\r\\nAccept-Language: en,zh-CN;q=0.9,zh;q=0.8,it;q=0.7\\r\\nCache-Control: no-cache\\r\\nConnection: keep-alive\\r\\nContent-Type: application/json\\r\\nOrigin: chrome-extension://fhbjgbiflinjbdggehcddcbncdddomop\\r\\nPostman-Token: 1938af14-db7a-26ad-541d-76aca94ab796\\r\\nSec-Ch-Ua: \\\"Chromium\\\";v=\\\"104\\\", \\\" Not A;Brand\\\";v=\\\"99\\\", \\\"Google Chrome\\\";v=\\\"104\\\"\\r\\nSec-Ch-Ua-Mobile: ?0\\r\\nSec-Ch-Ua-Platform: \\\"macOS\\\"\\r\\nSec-Fetch-Dest: empty\\r\\nSec-Fetch-Mode: cors\\r\\nSec-Fetch-Site: none\\r\\n\\r\\n{\\n\\t\\\"branch\\\": \\\"fix-bugs\\\", \\n\\t\\\"popIndex\\\": \\\"3\\\", \\n\\t\\\"location\\\": \\\"aws94\\\", \\n\\t\\\"hostname\\\": \\\"aws94-jaj\\\"\\n}\"}"

	var remoteReq RemoteRequest
	if err := json.Unmarshal([]byte(data), &remoteReq); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(remoteReq)
	}
}
