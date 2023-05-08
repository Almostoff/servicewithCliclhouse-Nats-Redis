package database

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"net/http"
	"taskFive/server/config"
	"taskFive/server/internal/entity"
	e "taskFive/server/lib/err"
	"time"
)

type Database struct {
	db *sql.DB
}

var m int = 1

func NewDb() *Database {
	return &Database{}
}

func InitDBConn(cfg config.Config) (*Database, error) {
	dbConn := Database{}
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Addr, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBname)
	dbConn.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return &Database{}, err
	}

	defer func() { err = e.WrapIfErr("can`t init database connection: ", err) }()

	return &dbConn, err
}

func (d Database) CreateItem(item *entity.Item, cId int) (*entity.Item, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO items (campaign_id, name, description, priority, removed, created_at)
          VALUES ($1, $2, $3, $4, $5, $6)
          RETURNING id`

	err = tx.QueryRow(q, cId, item.Name, item.Description, m, item.Removed, time.Now().Format(http.TimeFormat)).Scan(&item.ID)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(item)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	m++

	defer func() { err = e.WrapIfErr("can`t create new note in table: ", err) }()

	return d.GetItemById(item.ID, cId)
}

func (d Database) PatchItem(item *entity.Item, cId, iId int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT * FROM items WHERE id=$1 AND campaign_id=$2 FOR UPDATE", iId, cId)
	data := new(entity.Item)
	err = row.Scan(&data.ID, &data.CampaignID, &data.Name, &data.Description, &data.Priority, &data.Removed, &data.CreatedAt)
	if err != nil {
		return err
	}

	item.CreatedAt = time.Now()

	validate := validator.New()
	err = validate.Struct(item)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE items SET name = $1, description = $2, priority = $3, removed = $4, created_at = $5 "+
		"WHERE id = $6 AND campaign_id = $7", data.Name, data.Description, data.Priority,
		data.Removed, data.CreatedAt.Format(http.TimeFormat), iId, cId)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	err = tx.Commit()
	if err != nil {
		return err
	}

	defer func() { err = e.WrapIfErr("can`t update item: ", err) }()

	return nil
}

func (d Database) GetItemById(id, cId int) (*entity.Item, error) {
	row := d.db.QueryRow("SELECT * FROM items WHERE id=$1 AND campaign_id = $2", id, cId)
	data := new(entity.Item)
	err := row.Scan(&data.ID, &data.CampaignID,
		&data.Name, &data.Description,
		&data.Priority, &data.Removed, &data.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, err
		}
		return nil, err
	}

	defer func() { err = e.WrapIfErr("can`t get item by id: ", err) }()

	return data, err
}

func (d Database) DeleteItem(iId, cId int) (*entity.Item, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("SELECT * FROM items WHERE id=$1 AND campaign_id=$2 FOR UPDATE", iId, cId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	time.Sleep(time.Second * 10)

	_, err = tx.Exec("UPDATE items SET removed = $1 WHERE id = $2 AND campaign_id = $3", true, iId, cId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	row := tx.QueryRow("SELECT * FROM items WHERE id=$1 AND campaign_id=$2", iId, cId)
	data := new(entity.Item)
	err = row.Scan(&data.ID, &data.CampaignID,
		&data.Name, &data.Description,
		&data.Priority, &data.Removed, &data.CreatedAt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM items WHERE id = $1", iId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	defer func() { err = e.WrapIfErr("can`t delete item: ", err) }()

	return data, nil
}

func (d Database) GetAll() ([]*entity.Item, error) {
	rows, err := d.db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	itemList := []*entity.Item{}
	for rows.Next() {
		item := &entity.Item{}
		err := rows.Scan(&item.ID, &item.CampaignID, &item.Name,
			&item.Description, &item.Priority, &item.Removed, &item.CreatedAt)
		if err != nil {
			return itemList, err
		}
		itemList = append(itemList, item)
	}

	defer func() { err = e.WrapIfErr("can`t get item list: ", err) }()

	return itemList, err
}

func (d Database) RowsCount() int {
	var count int

	_ = d.db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)

	return count
}
