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

func messageReceived(client *MQTT.Client, msg MQTT.Message) {
	topics := strings.Split(msg.Topic(), "/")
	msgFrom := topics[len(topics)-1]
	fmt.Print(msgFrom + ": " + string(msg.Payload()))
}

func main() {
	stdin := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().Unix())

	server := flag.String("server", "tcp://iot.eclipse.org:1883", "The MQTT server to connect to")
	room := flag.String("room", "gochat", "The chat room to enter. default 'gochat'")
	name := flag.String("name", "user"+strconv.Itoa(rand.Intn(1000)), "Username to be displayed")
	flag.Parse()

	subTopic := strings.Join([]string{"/gochat/", *room, "/+"}, "")
	pubTopic := strings.Join([]string{"/gochat/", *room, "/", *name}, "")

	opts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*name).SetCleanSession(true)

	opts.OnConnect = func(c *MQTT.Client) {
		if token := c.Subscribe(subTopic, 1, messageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Connected as %s to %s\n", *name, *server)

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		if token := client.Publish(pubTopic, 1, false, message); token.Wait() && token.Error() != nil {
			fmt.Println("Failed to send message")
		}
	}
}
