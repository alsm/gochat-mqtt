gochat-mqtt
===========

A simple chat program to demonstrate [MQTT](http://mqtt.org) with Go, using the [Eclipse Paho](http://eclipse.org/paho) Go client

*build*
(requires Go, obviously...)
Set `GOPATH` and `GOBIN` variables
```
$ go get git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git
$ go install mqttchat.go
```
*run*
```
$ ./bin/mqttchat --help
Usage of ./bin/mqttchat:
  -name="user201": Username to be displayed
  -room="gochat": The chat room to enter. default 'gochat'
  -server="tcp://iot.eclipse.org:1883": The MQTT server to connect to
```
