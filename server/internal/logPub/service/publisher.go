package service

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"taskFive/server/internal/entity"
	e "taskFive/server/lib/err"
)

const (
	NatsStrUrl = "localhost:4224"
	clusterId  = "test-cluster"
	clientId   = "test-publisher"
	chann      = "natschannel"
)

type StanClient struct {
	sc stan.Conn
}

func CreateStan() *StanClient {
	stan, err := stan.Connect(clusterId, clientId, stan.NatsURL(NatsStrUrl))
	if err != nil {
		log.Print("can`t create connect to nats", err)
		panic(err)
	}
	log.Println("connection to nats successful")
	return &StanClient{sc: stan}
}

func (sC *StanClient) PublishLogMessage(data []byte) error {
	a := &entity.LogData{}
	err := json.Unmarshal(data, a)
	if err != nil {
		return err
	}

	err = sC.sc.Publish(chann, data)
	if err != nil {
		return err
	}

	defer func() { err = e.WrapIfErr("can`t publish log message: ", err) }()

	log.Println("send message to nats")
	return nil
}
