package accessory

import (
	"github.com/boisjacques/hc/log"
	"github.com/boisjacques/hc/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	log.Debug.Printf("Split topic:\t%v", topic)
	stringified := string(msg.Payload())
	value, err := strconv.ParseFloat(stringified, 32)
	if err != nil {
		log.Debug.Fatalf("Error handling Message %v\nMessage String: %v\n", msg.Payload(), stringified)
	}
	switch topic {
	case "temperature":
		a.handleTemperature(value)
	case "humidity":
		a.handleHumidity(value)
	}
}

func (a *MqttAccessory) handleTemperature(temperature float64) {
	a.TempSensor.CurrentTemperature.SetValue(temperature)
	log.Info.Printf("Temperature set to %f\n", temperature)
}

func (a *MqttAccessory) handleHumidity(humidity float64) {
	a.HumiditySensor.CurrentRelativeHumidity.SetValue(humidity)
	log.Info.Printf("Humidity set to %f\n", humidity)
}

func NewMqttAccessory(info Info) *MqttAccessory {
	acc := MqttAccessory{}
	acc.Accessory = New(info, info.DeviceType)
	if info.DeviceType == 5 {
		acc.Light = service.NewColoredLightbulb()
		acc.Light.Hue.SetValue(0)
		acc.Light.Saturation.SetValue(0)
		acc.Light.Brightness.SetValue(0)
		acc.AddService(acc.Light.Service)
		log.Debug.Printf("Topic:\t%v\tService:\t%v\n", acc.Type, acc.Light.Type)
	} else if info.DeviceType == 9 {
		acc.TempSensor = service.NewTemperatureSensor()
		acc.TempSensor.CurrentTemperature.SetValue(20)
		acc.TempSensor.CurrentTemperature.SetMinValue(-50)
		acc.TempSensor.CurrentTemperature.SetMaxValue(50)
		acc.TempSensor.CurrentTemperature.SetStepValue(0.1)
		acc.AddService(acc.TempSensor.Service)
		log.Debug.Printf("Topic:\t%v\tService:\t%v\n", acc.Type, acc.TempSensor.Type)
		acc.HumiditySensor = service.NewHumiditySensor()
		acc.HumiditySensor.CurrentRelativeHumidity.SetValue(50)
		acc.HumiditySensor.CurrentRelativeHumidity.SetMinValue(0)
		acc.HumiditySensor.CurrentRelativeHumidity.SetMaxValue(100)
		acc.HumiditySensor.CurrentRelativeHumidity.SetStepValue(0.1)
		acc.Accessory.AddService(acc.HumiditySensor.Service)
	}
	return &acc
}
