package sub

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"taskFive/server/clickhouse"
	"taskFive/server/config"
	"taskFive/server/internal/entity"
	e "taskFive/server/lib/err"
)

type StanSub struct {
	sc stan.Conn
	cH clickhouse.House
}

func CreateSub(cfg config.Config, cl clickhouse.House) *StanSub {
	stan, err := stan.Connect(cfg.Nats.ClusterID, cfg.Nats.ClientID, stan.NatsURL(cfg.Nats.Url))
	if err != nil {
		log.Println(e.WrapIfErr("can`t do command: save page", err))
		return nil
	}
	log.Println("connection nats subscriber successful")
	return &StanSub{
		sc: stan,
		cH: cl,
	}
}

func (sSub *StanSub) SubscribeToChannel(channel string, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	sChan := make(chan *stan.Msg, 10)
	go func() {
		var messages []*stan.Msg
		for msg := range sChan {
			messages = append(messages, msg)
			if len(messages) >= 2 {
				err := sSub.sendBatchToClickHouse(messages)
				if err != nil {
					log.Printf("error sending messages to ClickHouse: %v", err)
				}
				messages = nil
			}
		}
	}()
	sub, err := sSub.sc.Subscribe(channel, func(msg *stan.Msg) {
		log.Println("RECEIVED A NEW MESSAGE FROM NATS -")
		fmt.Println("Received a message: ", string(msg.Data))
		sChan <- msg
	}, opts...)
	if err != nil {
		return nil, err
	}
	defer func() { err = e.WrapIfErr("can`t subs to new messages from nats channel: ", err) }()
	return sub, err
}

func (sSub *StanSub) handlerMsg(msg *stan.Msg) {
	fmt.Println("Received a message: ", string(msg.Data))
	msg.Ack()
}

func (sSub *StanSub) sendBatchToClickHouse(messages []*stan.Msg) error {
	var values []*entity.LogData
	for _, msg := range messages {
		var logData *entity.LogData
		err := json.Unmarshal(msg.Data, &logData)
		if err != nil {
			log.Printf("Error unmarshalling message %s: %v", string(msg.Data), err)
			continue
		}
		values = append(values, logData)
	}

	err := sSub.cH.InsertLog(values)
	if err != nil {
		return err
	}

	defer func() { err = e.WrapIfErr("can`t send batch to click house: ", err) }()

	return nil
}
