package main

import (
  "log"
	"encoding/json"
	"time"
	"net/http"
	"fmt"
	"github.com/andrewn/realtime-analytics"
)

func main() {
	log.Printf("main")

	broker := realtime.NewServer()
	analytics := realtime.NewAnalytics()
	
	log.Printf("pre-go")
	
	go func() {
		for {
			data := analytics.GetData()
			rows := data.Rows
			
			jsonified, err := json.Marshal(rows)
			
		  if err != nil {
				fmt.Println(err)
			}
			
			log.Printf("Sending %o rows", data.TotalResults)

			broker.Notifier <- []byte(jsonified)
			
			time.Sleep(time.Second * 10)
		}
	}()


	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", broker))

}