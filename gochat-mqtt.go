package main

import (
	"bufio"
	"flag"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	server := flag.String("server", "tcp://iot.eclipse.org:1883", "The MQTT server to connect to")
	room := flag.String("room", "gochat", "The chat room to enter. default 'gochat'")
	name := flag.String("name", "user"+strconv.Itoa(rand.Intn(1000)), "Username to be displayed")
	flag.Parse()

	opts := MQTT.NewClientOptions().AddBroker(*server).SetClientId(*name).SetCleanSession(true)
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected as %s to %s\n", *name, *server)
	}

	sub_topic := strings.Join([]string{"/gochat/", *room, "/+"}, "")
	if token := client.Subscribe(sub_topic, 1,
		func(client *MQTT.MqttClient, msg MQTT.Message) {
			msg_from := strings.Split(msg.Topic(), "/")[3]
			fmt.Println(msg_from + ": " + string(msg.Payload()))
		}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	pub_topic := strings.Join([]string{"/gochat/", *room, "/", *name}, "")

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		client.Publish(pub_topic, 1, false, strings.TrimSpace(message)).Wait()
	}
}
