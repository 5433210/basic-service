package main

import (
	"context"
	"flag"
	"net"
	"os"

	"github.com/eclipse/paho.golang/paho"

	"wailik.com/internal/pkg/log"
)

func main() {
	stdout := log.New(log.OptOutputPaths([]string{"stdout"}))
	defer stdout.Flush()
	stderr := log.New(log.OptErrorOutputPaths([]string{"stderr"}))
	defer stderr.Flush()

	hostname, _ := os.Hostname()

	server := flag.String("server", "127.0.0.1:1883", "The full URL of the MQTT server to connect to")
	topic := flag.String("topic", hostname, "Topic to publish the messages on")
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	clientid := flag.String("clientid", "", "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	data := flag.String("data", "", "")
	flag.Parse()

	conn, err := net.Dial("tcp", *server)
	if err != nil {
		stderr.Fatalf("Failed to connect to %s: %s", *server, err)

		return
	}

	c := paho.NewClient(paho.ClientConfig{
		Conn: conn,
	})

	cp := &paho.Connect{
		KeepAlive:  30,
		ClientID:   *clientid,
		CleanStart: true,
		Username:   *username,
		Password:   []byte(*password),
	}

	if *username != "" {
		cp.UsernameFlag = true
	}
	if *password != "" {
		cp.PasswordFlag = true
	}

	ca, err := c.Connect(context.Background(), cp)
	if err != nil {
	}
	if ca.ReasonCode != 0 {
		stderr.Fatalf("Failed to connect to %s : %d - %s",
			*server, ca.ReasonCode, ca.Properties.ReasonString)

		return
	}

	if _, err = c.Publish(context.Background(), &paho.Publish{
		Topic:   *topic,
		QoS:     byte(*qos),
		Retain:  *retained,
		Payload: []byte(*data),
	}); err != nil {
		stderr.Fatalf("Failed to publish message:", err)

		return
	}

	c.Conn.Close()

	stdout.Infof("sent message successfully:%+v", *data)
}
