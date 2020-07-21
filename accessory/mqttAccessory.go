package accessory

import (
	"github.com/boisjacques/hc/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
)

type MqttAccessory struct {
	*Accessory
	kind           string
	rxTopics       []string
	txTopics       []string
	TempSensor     *service.TemperatureSensor
	HumiditySensor *service.HumiditySensor
	Light          *service.ColoredLightbulb
}

func (a *MqttAccessory) HandleMessage(msg mqtt.Message) {
	topic := msg.Topic()
	topic = strings.Split(topic, "/")[2]
	stringified := string(msg.Payload())
	if value, err := strconv.ParseFloat(stringified, 32); err != nil {
		switch topic {
		case "temperature":
			a.handleTemperature(value)
		case "humidity":
			a.handleHumidity(value)
		}
	}
}

func (a *MqttAccessory) handleTemperature(temperature float64) {
	a.TempSensor.CurrentTemperature.SetValue(temperature)
	log.Printf("Temperature set to %f\n", temperature)
}

func (a *MqttAccessory) handleHumidity(humidity float64) {
	a.TempSensor.CurrentTemperature.SetValue(humidity)
	log.Printf("Temperature set to %f\n", humidity)
}

func NewMqttAccessory(info Info) *MqttAccessory {
	acc := MqttAccessory{}
	typ, _ := strconv.Atoi(info.Model)
	acc.Accessory = New(info, AccessoryType(typ))
	if typ == 5 {
		acc.Light = service.NewColoredLightbulb()
		acc.Light.Hue.SetValue(0)
		acc.Light.Saturation.SetValue(0)
		acc.Light.Brightness.SetValue(0)
		acc.AddService(acc.Light.Service)
	} else if typ == 9 {
		acc.TempSensor = service.NewTemperatureSensor()
		acc.TempSensor.CurrentTemperature.SetValue(20)
		acc.TempSensor.CurrentTemperature.SetMinValue(-50)
		acc.TempSensor.CurrentTemperature.SetMaxValue(50)
		acc.TempSensor.CurrentTemperature.SetStepValue(0.1)
		acc.AddService(acc.TempSensor.Service)
		/*
			acc.HumiditySensor = service.NewHumiditySensor()
			acc.HumiditySensor.CurrentRelativeHumidity.SetValue(50)
			acc.HumiditySensor.CurrentRelativeHumidity.SetMinValue(0)
			acc.HumiditySensor.CurrentRelativeHumidity.SetMaxValue(100)
			acc.HumiditySensor.CurrentRelativeHumidity.SetStepValue(0.1)
			acc.Accessory.AddService(acc.HumiditySensor.Service)
		*/
	}
	return &acc
}
