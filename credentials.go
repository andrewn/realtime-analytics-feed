package realtime

import (
  "fmt"
  "io/ioutil"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  analytics "google.golang.org/api/analytics/v3"
)

// https://godoc.org/google.golang.org/api/analytics/v3
func CreateAnalyticsCall() (realtimeServiceGetCall *analytics.DataRealtimeGetCall) {
  scope := "https://www.googleapis.com/auth/analytics.readonly"

  raw, err := ioutil.ReadFile("./credentials.json")
  
  if err != nil {
    fmt.Println(err)
  }
  
  conf, err := google.JWTConfigFromJSON(raw, scope)

  if err != nil {
    fmt.Println(err)
  }  

  client := conf.Client(oauth2.NoContext)
  
  analyticsService, err := analytics.New(client)
  
  if err != nil {
    fmt.Println(err)
  }
  
  realtimeService := analytics.NewDataRealtimeService(analyticsService)
  
  realtimeServiceGetCall = realtimeService.Get("ga:25221044", "rt:activeUsers")
  realtimeServiceGetCall.Dimensions("rt:pagePath")
  
  return
}