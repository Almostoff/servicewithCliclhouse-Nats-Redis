package clickhouse

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"log"
	"taskFive/server/config"
	"taskFive/server/internal/entity"
	e "taskFive/server/lib/err"
)

type House struct {
	ch *sql.DB
}

func NewClickHouseDB(cfg config.Config) (*House, error) {
	clConn := House{}
	var err error
	cHouseInfo := fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s",
		cfg.CHouse.Host, cfg.CHouse.Port, cfg.CHouse.Username, cfg.CHouse.Password, cfg.CHouse.Database)
	clConn.ch, err = sql.Open("clickhouse", cHouseInfo)
	if err != nil {
		return nil, err
	}
	log.Println("connect to clickHouse database successful")
	defer func() { err = e.WrapIfErr("can`t connect to clickhouse", err) }()
	return &clConn, nil
}

func (c *House) InsertLog(data []*entity.LogData) error {
	tx, err := c.ch.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO log (Id, CampaignId, Name, Description, Priority, Removed, EventTime) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, logData := range data {
		_, err = stmt.Exec(logData.ID, logData.CampaignId, logData.Name, logData.Description, logData.Priority, logData.Removed, logData.EventTime)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("log batch send to clickhouse database successful")

	defer func() { err = e.WrapIfErr("can`t insert log", err) }()

	return nil
}
