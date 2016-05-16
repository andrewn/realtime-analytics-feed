package realtime

import (
  "log"
  "fmt"
  // "time"
  analytics "google.golang.org/api/analytics/v3"
)

type Analytics struct {
  get *analytics.DataRealtimeGetCall
}

func NewAnalytics() (analytics *Analytics) {
  analytics = &Analytics{
    get: CreateAnalyticsCall(),
  }
  
  return
}

func (analytics *Analytics) GetData() (response *analytics.RealtimeData) {
  response, err := analytics.get.Do()
  
  if err != nil {
    fmt.Println(err)
  }
  
  log.Printf("Results %o", response.TotalResults)
  log.Printf("Rows %o", response.Rows)
  
  return
}