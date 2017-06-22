package main

import (
	"flag"
	"./flesapi"
	"log"
	"encoding/json"
	"os"
)

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func main() {
	tokenFlag := flag.String("token", "", "Flespi auth token")
	limitCountFlag := flag.Int("limit_count", 10000, "Limit count")
	//limitSizeFlag := flag.Int("limit_size", 1000, "Limit size")
	timeoutFlag := flag.Int("timeout", 10, "timeout between request")
	deleteFlag := flag.Bool("delete", false, "Delete message")

	channelFlag := flag.Int("channel", 0, "Channeld id")
	flag.Parse()


	client := flesapi.NewClient(*tokenFlag)

	args := flesapi.MessageArgs{
		Curr_key: 0,
		Limit_count: *limitCountFlag,
		//Limit_size: *limitSizeFlag,
		Timeout: *timeoutFlag,
		Delete: *deleteFlag}

	f, err := os.Create("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(f)

	for {
		response, err := client.GetChannelMessages(*channelFlag, args)

		if err != nil {
			log.Fatal(err)
		}

		for _, res := range response.Result {
			encoder.Encode(res)
		}

		if response.Next_key == 0 {
			break
		}

		args.Curr_key = max(args.Curr_key, response.Next_key)
	}
	f.Close()
}
