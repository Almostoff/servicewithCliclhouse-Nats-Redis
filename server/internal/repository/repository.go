package repository

import (
	"encoding/json"
	"log"
	"strconv"
	"taskFive/server/internal/database"
	"taskFive/server/internal/entity"
	"taskFive/server/internal/logPub/service"
	"taskFive/server/internal/redis"
	e "taskFive/server/lib/err"
	"time"
)

type Repository struct {
	db  database.Database
	red redis.Redis
	n   service.StanClient
}

func InitStore(db database.Database, red redis.Redis, nats service.StanClient) *Repository {
	Repository := Repository{
		db:  db,
		red: red,
		n:   nats,
	}
	return &Repository
}

func serializeLogData(data *entity.Item) *entity.LogData {
	var b uint8
	if data.Removed == true {
		b = 1
	} else {
		b = 0
	}
	return &entity.LogData{
		ID:          data.ID,
		CampaignId:  data.CampaignID,
		Name:        data.Name,
		Description: data.Description,
		Priority:    data.Priority,
		Removed:     b,
		EventTime:   data.CreatedAt,
	}
}

func (r Repository) CreateItem(item *entity.Item, campaignId int) (*entity.Item, error) {
	data, err := r.db.CreateItem(item, campaignId)
	if err != nil {
		return nil, e.WrapIfErr("can`t create new item: ", err)
	}

	l := serializeLogData(data)
	logValue, err := json.Marshal(&l)
	if err != nil {
		return nil, e.WrapIfErr("can`t serialize log data", err)
	}

	if err = r.n.PublishLogMessage(logValue); err != nil {
		return nil, e.WrapIfErr("can`t send message to nats server", err)
	}

	value, err := json.Marshal(&data)
	if err != nil {
		return nil, e.WrapIfErr("can`t serialize data", err)
	}

	if err = r.red.SetValue(strconv.Itoa(data.ID), value, time.Minute); err != nil {
		return nil, e.WrapIfErr("can`t save data to Redis", err)
	}

	return data, err
}

func (r Repository) PatchItem(item *entity.Item, campaignId, iId int) (*entity.Item, error) {
	_, err := r.db.GetItemById(iId, campaignId)
	if err != nil {
		return nil, err
	}
	err = r.db.PatchItem(item, campaignId, iId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = r.red.DeleteValue(strconv.Itoa(iId))
	if err != nil {
		log.Print("can`t delete from Redis: ", err)
		return nil, err
	}

	data, err := r.db.GetItemById(iId, campaignId)
	value, err := json.Marshal(&data)
	if err != nil {
		log.Fatal("can`t serialize data", err)
	}

	l := serializeLogData(data)
	logValue, err := json.Marshal(&l)
	if err != nil {
		log.Fatal("can`t serialize data", err)
	}

	err = r.n.PublishLogMessage(logValue)

	_ = r.red.SetValue(strconv.Itoa(data.ID), value, time.Minute)

	defer func() { err = e.WrapIfErr("can`t patch item: ", err) }()

	return data, err
}

func (r Repository) DeleteItem(iId, cId int) (*entity.DelR, error) {
	_, err := r.db.GetItemById(iId, cId)
	if err != nil {
		return nil, err
	}
	data, err := r.db.DeleteItem(iId, cId)

	l := serializeLogData(data)
	logValue, err := json.Marshal(&l)
	if err != nil {
		log.Fatal("can`t serialize data", err)
	}

	err = r.n.PublishLogMessage(logValue)

	err = r.red.DeleteValue(strconv.Itoa(iId))
	if err != nil {
		log.Print("can`t delete from Redis: ", err)
		return nil, err
	}

	dR := &entity.DelR{
		ID:         data.ID,
		CampaignId: data.CampaignID,
		Removed:    data.Removed,
	}

	defer func() { err = e.WrapIfErr("can`t delete item: ", err) }()

	return dR, err
}

func (r Repository) GetAll() ([]*entity.Item, error) {
	var itemList []*entity.Item
	data, _, _ := r.red.GetAllValue()
	if data != nil && len(data) == r.db.RowsCount() {
		log.Println("read from redis")
		for _, v := range data {
			var item *entity.Item
			err := json.Unmarshal([]byte(v.(string)), &item)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			itemList = append(itemList, item)
		}
		return itemList, nil
	} else {
		data, err := r.db.GetAll()
		if err != nil {
			return nil, err
		}
		log.Println("read from database")
		pairs := make([]interface{}, 0, len(data)*2)
		for _, item := range data {
			key := strconv.Itoa(item.ID)
			value, err := json.Marshal(item)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			pairs = append(pairs, key, string(value))
		}
		if err := r.red.SetManyValue(pairs); err != nil {
			log.Print(err)
		}

		defer func() { err = e.WrapIfErr("can`t take all items: ", err) }()

		return data, nil
	}
}
