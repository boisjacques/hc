package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Bridge struct {
	uri      string
	username string
	password string
}

func NewMQTTBridge() *Bridge {
	bridge := Bridge{}
	return &bridge
}

func (mb *Bridge) Register(username string, password string, clientId string, uri string) mqtt.Client {
	opts := mb.createClientOptions(username, password, clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func (mb *Bridge) createClientOptions(username string, password string, clientId string, uri string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}
