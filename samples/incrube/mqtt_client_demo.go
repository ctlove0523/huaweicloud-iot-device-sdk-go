package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {
	ops := &mqtt.ClientOptions{}
	ops.AddBroker("tcp://localhost:8883")
	ops.SetClientID("test-republish")
	ops.SetUsername("test")
	ops.SetPassword("hello")
	ops.SetKeepAlive(5 * time.Second)
	//ops.SetPingTimeout(2 * time.Second)

	client := mqtt.NewClient(ops)
	token := client.Connect()
	token.Wait()

	client.Subscribe("test", 1, func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("get message %s\n", string(message.Payload()))
		message.Ack()
	})

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(time.Now().String())
		fmt.Println(client.IsConnected())
	}
	ch := make(chan struct{})
	<-ch

}
