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

	opts := MQTT.NewClientOptions().SetBroker(*server).SetClientId(*name).SetCleanSession(true).SetTraceLevel(MQTT.Off)
	client := MQTT.NewClient(opts)
	_, err := client.Start()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Connected as %s to %s\n", *name, *server)
	}

	client.StartSubscription(func(msg MQTT.Message) {
		msg_from := strings.Split(msg.Topic(), "/")[3]
		fmt.Println(msg_from + ": " + string(msg.Payload()))
	}, "/gochat/"+*room+"/+", MQTT.QOS_ONE)

	pub_topic := strings.Join([]string{"/gochat/", *room, "/", *name}, "")

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		r := client.Publish(MQTT.QOS_ONE, pub_topic, strings.TrimSpace(message))
		<-r
	}
}
