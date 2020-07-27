package main

import (
	"flag"
	"github.com/boisjacques/hc"
	"github.com/boisjacques/hc/accessory"
	"github.com/boisjacques/hc/log"
	"github.com/boisjacques/hc/mqtt"
	"sync"
)

func main() {
	username := flag.String("username", "mqttuser", "mqtt user")
	password := flag.String("password", "mqttpassword", "mqtt password")
	uri := flag.String("uri", "tcp://192.168.2.250:1883", "mqtt broker uri")
	cid := flag.String("cid", "test-client", "mqtt client id")
	flag.Parse()
	log.Debug.Enable()

	accessories := make([]accessory.MqttAccessory, 0)
	bridge := *mqtt.NewMQTTBridge()
	var wg sync.WaitGroup

	// TODO: Write json config with accessories
	infos := []accessory.Info{
		{
			Name:         "projectada",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   accessory.TypeLightbulb,
			Topics: []string{
				"home/manu/light/on",
				"home/manu/light/hue",
				"home/manu/light/saturation",
				"home/manu/light/brightness",
			},
		},
		{
			Name:         "Manucave",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   accessory.TypeThermostat,
			Topics: []string{
				"home/manu/temperature",
				"home/manu/humidity",
			},
		},
		{
			Name:         "Mancave",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   accessory.TypeThermostat,
			Topics: []string{
				"home/mancave/temperature",
				"home/mancave/humidity",
			},
		},
		{
			Name:         "Basement",
			Manufacturer: "HoChiMinh Flowerpower Enterprises",
			DeviceType:   accessory.TypeThermostat,
			Topics: []string{
				"home/basement/temperature",
				"home/basement/humidity",
			},
		},
	}

	client := bridge.Register(*username, *password, *cid, *uri)
	dispatcher := mqtt.NewDispatcher(client)
	publishChannel := make(chan string)

	for _, info := range infos {
		acc := accessory.NewMqttAccessory(info, publishChannel)
		accessories = append(accessories, *acc)
		if info.DeviceType != accessory.TypeLightbulb {
			for _, topic := range info.Topics {
				dispatcher.Subscribe(topic, *acc)
			}
		}
	}
	dispatcher.Publish(publishChannel)

	hkbridge := accessory.NewBridge(accessory.Info{
		Name:         "Bridge",
		SerialNumber: "1312",
	})
	t, err := hc.NewIPTransport(hc.Config{Pin: "11223344"}, hkbridge.Accessory, accessories[0].Accessory)
	if err != nil {
		log.Debug.Fatalln(err)
	}
	go dispatcher.Listen()

	hc.OnTermination(func() {
		<-t.Stop()
	})
	wg.Add(1)
	go t.Start(&wg)
	wg.Wait()
}
