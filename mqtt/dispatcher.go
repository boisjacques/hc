package mqtt

import (
	"github.com/boisjacques/hc/accessory"
	"github.com/boisjacques/hc/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Dispatcher struct {
	client     mqtt.Client
	deviceMap  map[string]accessory.MqttAccessory
	rxChannels []chan []byte
	txChannels []chan []byte
	rx         chan mqtt.Message
}

func NewDispatcher(client mqtt.Client) *Dispatcher {
	return &Dispatcher{
		client:    client,
		deviceMap: make(map[string]accessory.MqttAccessory),
		rx:        make(chan mqtt.Message, 100),
	}
}

func (d *Dispatcher) AddTopic(topic string, acc accessory.MqttAccessory) {
	d.deviceMap[topic] = acc
	go func() {
		d.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			log.Debug.Println("Waiting for message...")
			d.rx <- msg
		})
	}()
}

func (d *Dispatcher) Listen() {
	for {
		msg := <-d.rx
		log.Debug.Printf("Topic:\t%v\nPayload:\t%v", msg.Topic(), msg.Payload())
		mqttAccessory := d.deviceMap[msg.Topic()]
		mqttAccessory.HandleMessage(msg)
	}
}
