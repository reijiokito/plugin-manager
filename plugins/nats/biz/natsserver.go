package biz

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/nats-io/nats.go"
	"os"
)

type NatsServer struct {
	Broker *nats.Conn
	Logger hclog.Logger
}

type ClientStatus string

const (
	CLIENT_CONNECTED ClientStatus = "CONNECTED"
	CLIENT_CLOSED    ClientStatus = "CLOSED"
)

type ClientInfo struct {
	Name   string
	Status ClientStatus
	Config plugin.Client
}

func (n *NatsServer) Close(clientName string) string {
	n.Logger.Debug(" received message closed successfully" + clientName)

	return "Closed"
}

func (n *NatsServer) SendData(clientName string, data string) string {
	n.Logger.Debug("message from NatsConnect.SendData")

	n.Logger.Debug(fmt.Sprintf("nats: connected = %v, closed = %v", n.Broker.IsConnected(), n.Broker.IsClosed()))
	err := n.Broker.Publish(clientName, []byte(data))
	if err != nil {
		return err.Error()
	}

	n.Logger.Debug(" received message successfully")

	//select {}
	return "Khong co loi"
}

func (n *NatsServer) Subscript(subject string) string {
	if _, err := n.Broker.Subscribe(subject, func(m *nats.Msg) {
		n.Logger.Debug("Received a message: %s\n", string(m.Data))
	}); err != nil {
		n.Logger.Debug("ERR: %s\n", err.Error())
		return "SCRIPTION ERR: " + subject
	}
	select {}

	return "END SUBSCRIPTION"
}

func (n *NatsServer) Connect(name string, config plugin.Client) string {
	file, err := os.Create("client_" + name + ".gob")
	if err != nil {
		fmt.Println("LOI O DAY: " + err.Error())
		return "Connect error: Sent from plugin"
	}
	defer file.Close()

	// create a new encoder object
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("HERE: " + err.Error())
		return "Setup config error: Sent from plugin"
	}

	return "Client Connected. Sent from Server"
	return "Khong co loi"
}

func getClient(name string) (ClientInfo, error) {
	fi, err := os.Open("client_" + name + ".gob")
	if err != nil {
		return ClientInfo{}, errors.New("Get Client Err: " + err.Error())
	}
	defer fi.Close()

	decoder := gob.NewDecoder(fi)
	var client ClientInfo
	err = decoder.Decode(&client)
	if err != nil {
		fmt.Println("Decode err : ", err.Error())
		return ClientInfo{}, errors.New("Get Client Err: " + err.Error())
	}

	return client, nil

}
