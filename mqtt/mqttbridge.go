package mqtt

import (
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Bridge struct {

}

func (mb *Bridge) connect(username string, password string,clientId string, uri *url.URL) mqtt.Client {
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

func (mb *Bridge) createClientOptions(username string, password string, clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://192.168.2.252:1883")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func (mb *Bridge) listen(uri *url.URL, username string, password string, topic string, clientId string, c chan []byte) {
	client := mb.connect(username, password, clientId, uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Println(msg)
		c <- msg.Payload()
	})
}

func NewMQTTBridge(username string, password string, topic string, clientId string, c chan []byte) *Bridge {
	bridge := Bridge{}
	uri, err := url.Parse("tcp://192.168.2.252:1883")
	if err != nil {
		log.Fatal(err)
	}

	go bridge.listen(uri, username, password, topic, clientId, c)

	return &bridge
}

