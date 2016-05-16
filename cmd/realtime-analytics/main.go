package main

import (
  "log"
	"os"
	"encoding/json"
	"time"
	"net/http"
	"fmt"
	"github.com/andrewn/realtime-analytics"
)

func main() {
	log.Printf("Starting realtime-analytics main")
	
	port := os.Getenv("PORT")
	
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	broker := realtime.NewServer()
	analytics := realtime.NewAnalytics()
	
	
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

	listenOn := ":"+port
	
	log.Println("Listen on port", listenOn)

	log.Fatal("HTTP server error: ", http.ListenAndServe(listenOn, broker))

}