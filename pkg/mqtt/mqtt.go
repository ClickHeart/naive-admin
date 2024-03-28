package mqtt

import (
	"fmt"
	"math/rand"
	"naive-admin/pkg/config"
	"sync"
	"time"

	Mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	Mqtt.Client
}

func (client *MqttClient) MqttSubscribe(topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

func (client *MqttClient) MqttPublish(topic string, playload interface{}) {
	client.Publish(topic, 0, false, playload)
}

var messagePubHandler Mqtt.MessageHandler = func(client Mqtt.Client, msg Mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler Mqtt.OnConnectHandler = func(client Mqtt.Client) {
	fmt.Println("mqtt Connected")
}

var connectLostHandler Mqtt.ConnectionLostHandler = func(client Mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var Cli *MqttClient

var once sync.Once

func NewClient(cfg *config.Mqtt) *MqttClient {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientID := fmt.Sprintf("go-client-%d", r.Int())

	opts := Mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Second * 60)
	opts.SetPingTimeout(1 * time.Second)

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	if cfg.User != "" {
		opts.SetUsername(cfg.User)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	mqttClient := Mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &MqttClient{Client: mqttClient}
}

func Init() {
	singleton := func() {
		if config.Conf == nil {
			panic("config is nil!")
		}
		Cli = NewClient(config.Conf.Server.Mqtt)
	}
	once.Do(singleton)

}
