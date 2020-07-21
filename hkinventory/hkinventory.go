package main

import (
	"flag"
	"github.com/boisjacques/hc"
	"github.com/boisjacques/hc/accessory"
	"github.com/boisjacques/hc/mqtt"
	"log"
	"sync"
)

func main() {
	username := flag.String("username", "mqttuser", "mqtt user")
	password := flag.String("password", "mqttpassword", "mqtt password")
	uri := flag.String("uri", "tcp://192.168.2.250:1883", "mqtt broker uri")
	cid := flag.String("cid", "test-client", "mqtt client id")
	flag.Parse()

	accessories := make([]accessory.MqttAccessory, 0)
	bridge := *mqtt.NewMQTTBridge()
	var wg sync.WaitGroup

	// TODO: Write json config with accessories	}
	infos := []accessory.Info{
		{
			Name:         "Manucave",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   9,
			Topics: []string{
				"home/manu/temperature",
				"home/manu/humidity",
			},
		},
		{
			Name:         "Mancave",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   9,
			Topics: []string{
				"home/mancave/temperature",
				"home/mancave/humidity",
			},
		},
		{
			Name:         "Basement",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   9,
			Topics: []string{
				"home/basement/temperature",
				"home/basement/humidity",
			},
		},
	}

	client := bridge.Register(*username, *password, *cid, *uri)
	dispatcher := mqtt.NewDispatcher(client)

	for _, info := range infos {
		acc := accessory.NewMqttAccessory(info)
		accessories = append(accessories, *acc)
		for _, topic := range info.Topics {
			dispatcher.AddTopic(topic, *acc)
		}
	}

	hkbridge := accessory.NewBridge(accessory.Info{
		Name:             "Bridge",
		SerialNumber:     "1312",
	})
	t, err := hc.NewIPTransport(hc.Config{Pin: "11223344"}, hkbridge.Accessory, accessories[0].Accessory,
		accessories[1].Accessory,
		accessories[2].Accessory)
	if err != nil {
		log.Fatal(err)
	}
	go dispatcher.Listen()

	hc.OnTermination(func() {
		<-t.Stop()
	})
	wg.Add(1)
	go t.Start(&wg)
	wg.Wait()
}
