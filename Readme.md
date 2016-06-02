Polls the [Google Analytics Real Time Reporting API](https://developers.google.com/analytics/devguides/reporting/realtime/v3/) every 10s for active users on a website and pushes that data via Server Sent Events.

Configuration
---

Google Analytics API credentials are supplied by a JSON file that contains the private key for a Service Account.

0. Request beta access to the [Analytics Real Time Reporting API](https://developers.google.com/analytics/devguides/reporting/realtime/v3/).
1. [Create a Service Account](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount) for this application. The end result is a JSON file that's downloaded to you local machine.
2. In the Google Analytics Admin interface, add the `client_email` email address as a user to your Analytics account. This will give this application access to the API.
3. The JSON file can either be added as `credentials.json` in this directory or supplied as an environment variables called `CREDENTIALS_JSON` when starting the app.

Authorization
---

Set the AUTH_BASIC_USER and AUTH_BASIC_PASS env vars to require a username and password to access the event stream.

Deploy to Heroku
---

You can deploy this app to Heroku.

     heroku create
     git push -i heroku master
     heroku config:set GA_ID="ga:12345678"
     heroku config:set CREDENTIALS_JSON="`cat credentials.json`"

Visiting your heroku app URL should show the data being pushed.

Subscribe to the data source
---

This app is designed to use [Server Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) which makes it very easy to use in a web browser.

```js  
var source = new EventSource('http://your-app-name.herokuapp.com/')
source.onmessage = function (evt) {
  var data = JSON.parse(evt.data);
  console.log('New message: ', data);
}
```
