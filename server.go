package realtime

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool

	// auth config
	authConfig AuthConfig
}

type AuthConfig struct {
	BasicAuthUser string
	BasicAuthPass string
	BearerToken   string
}

func NewServer(auth AuthConfig) (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
		authConfig:     auth,
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func bearerTokenFromRequest(req *http.Request) (token string) {
	token = bearerTokenFromAuthorizationHeader(req)

	if token != "" {
		return
	}

	token = bearerTokenFromQueryString(req)

	return
}

func bearerTokenFromQueryString(req *http.Request) (token string) {
	query := req.URL.Query()

	if len(query["bearer"]) > 0 {
		token = query["bearer"][0]
	}

	return
}

func bearerTokenFromAuthorizationHeader(req *http.Request) (token string) {
	header := req.Header["Authorization"]

	if len(header) > 0 {
		return bearerTokenFromHeaderValue(header[0])
	}

	return
}

func bearerTokenFromHeaderValue(header string) (token string) {
	r, _ := regexp.Compile("Bearer +(\\S+)")
	matches := r.FindStringSubmatch(header)
	if len(matches) > 0 {
		token = matches[1]
	}
	return
}

func (broker *Broker) getCredentials() (user, pass, token string) {
	user = broker.authConfig.BasicAuthUser
	pass = broker.authConfig.BasicAuthPass
	token = broker.authConfig.BearerToken
	return
}

func (broker *Broker) hasAuthConfig() (hasAuthConfig bool) {
	user, pass, token := broker.getCredentials()
	hasBasicAuthConfig := user != "" && pass != ""
	hasBearerTokenConfig := token != ""
	hasAuthConfig = hasBasicAuthConfig || hasBearerTokenConfig
	return
}

func (broker *Broker) basicAuthMatch(suppliedUser, suppliedPass string, _ bool) (matches bool){
	matches = false
	user, pass, _ := broker.getCredentials()

	if suppliedUser == "" || suppliedPass == "" {
		return
	}

	if suppliedUser == user && suppliedPass == pass {
		matches = true
	}

	return
}

func (broker *Broker) bearerTokenMatch(suppliedToken string) (matches bool){
	matches = false
	_, _, token := broker.getCredentials()

	if suppliedToken == "" {
		return
	}

	if suppliedToken == token {
		matches = true
	}

	return
}

func (broker *Broker) IsAuthorised(req *http.Request) (isAuthorised bool) {
	isAuthorised = false

	fmt.Println("hasAuthConfig", broker.hasAuthConfig())

	if !broker.hasAuthConfig() {
		isAuthorised = true
	} else {
		if broker.basicAuthMatch(req.BasicAuth()) {
			fmt.Println("basic auth credentials match")
			isAuthorised = true
		} else if broker.bearerTokenMatch(bearerTokenFromRequest(req)) {
			fmt.Println("bearer token matches")
			isAuthorised = true
		}
		fmt.Println("isAuthorised", isAuthorised)
	}

	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	// Respond to CORS pre-flight request
	if req.Method == "OPTIONS" && req.Header.Get("Origin") != "" {
		fmt.Println("Is OPTIONS request with Origin")
		rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		return
	}

	if broker.hasAuthConfig() {
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// * precludes HTTP authorization so we explicitly allow Origin
	// supplied by client
	// rw.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Header.Get("Origin") != "" {
		rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	}

	isAuthorised := broker.IsAuthorised(req)

	fmt.Println("isAuthorised", isAuthorised)

	if !isAuthorised {
		rw.Header().Set("WWW-Authenticate", "Basic realm=\"realtime\"")
		http.Error(rw, "Incorrect credentials", http.StatusUnauthorized)
		return
	}

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

}

func (broker *Broker) listen() {
	var lastEvent []byte

	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))

			// Send lastEvent to newly connected client
			s <- lastEvent

		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:
			lastEvent = event
			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}
