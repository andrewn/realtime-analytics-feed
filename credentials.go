package realtime

import (
  "fmt"
  "os"
  "io/ioutil"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  analytics "google.golang.org/api/analytics/v3"
)

// https://godoc.org/google.golang.org/api/analytics/v3
func CreateAnalyticsCall() (realtimeServiceGetCall *analytics.DataRealtimeGetCall) {
  scope := "https://www.googleapis.com/auth/analytics.readonly"
  
  gaId := os.Getenv("GA_ID")
  
  if gaId == "" {
    fmt.Println("GA_ID not set")
    os.Exit(1)
  }
  
  raw, hasCredentials := GetCredentialsFromEnv()
  
  if hasCredentials == false {
    fmt.Println("Fetch credentials from file")
    raw = GetCredentialsFromFile()
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
  
  realtimeServiceGetCall = realtimeService.Get(gaId, "rt:activeUsers")
  realtimeServiceGetCall.Dimensions("rt:pagePath")
  
  return
}


func GetCredentialsFromEnv() (raw []byte, hasCredentials bool) {
  str := os.Getenv("CREDENTIALS_JSON")
  hasCredentials = false
  
  if str != "" {
    hasCredentials = true
    raw = []byte(str)
  }
  
  return
}

func GetCredentialsFromFile() (raw []byte) {
  raw, err := ioutil.ReadFile("./credentials.json")
  
  if err != nil {
    fmt.Println(err)
  }
  
  return
}
