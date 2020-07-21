package main

import (
	"flag"
	"github.com/boisjacques/hc"
	"github.com/boisjacques/hc/accessory"
	"github.com/boisjacques/hc/log"
	"github.com/boisjacques/hc/mqtt"
)

func main() {
	log.Debug.Enable()
	log.Debug.Println("Logging enabled")
	username := flag.String("username", "mqttuser", "mqtt user")
	password := flag.String("password", "mqttpassword", "mqtt password")
	uri := flag.String("uri", "tcp://192.168.2.250:1883", "mqtt broker uri")
	cid := flag.String("cid", "test-client", "mqtt client id")
	flag.Parse()
	bridge := mqtt.NewMQTTBridge()

	info := accessory.Info{
			Name:         "Manucave",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   9,
			Topics: []string{
				"home/manu/temperature",
				"home/manu/humidity",
			},
		}
	client := bridge.Register(*username, *password, *cid, *uri)
	dispatcher := mqtt.NewDispatcher(client)

	acc := accessory.NewMqttAccessory(info)
	for _, topic := range info.Topics {
		dispatcher.AddTopic(topic, *acc)
	}

	t, err := hc.NewIPTransport(hc.Config{Pin: "11223344"}, acc.Accessory)
	if err != nil {
		log.Debug.Fatalln(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
