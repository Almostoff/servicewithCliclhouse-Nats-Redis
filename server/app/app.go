package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"taskFive/server/clickhouse"
	"taskFive/server/config"
	"taskFive/server/internal/cases"
	"taskFive/server/internal/controller"
	"taskFive/server/internal/database"
	"taskFive/server/internal/logPub/service"
	"taskFive/server/internal/redis"
	"taskFive/server/internal/repository"
	"taskFive/server/sub"
	"time"
)

type App struct {
	cfg config.Config
}

func InitApp(cfg config.Config) *App {
	app := App{}
	app.cfg = cfg
	return &app
}

func (app *App) Run() {
	db, err := database.InitDBConn(app.cfg)
	if err != nil {
		log.Println("Error while connecting to database")
		panic(err)
	} else {
		log.Println("Connect to db successful")
	}

	cl, err := clickhouse.NewClickHouseDB(app.cfg)
	if err != nil {
		log.Println("Error while connecting to clickhouse database")
		panic(err)
	}

	red := redis.NewRedis(app.cfg)

	nats := service.CreateStan()

	repository := repository.InitStore(*db, *red, *nats)
	usecase := cases.NewUseCase(repository)
	r := chi.NewRouter()
	controller.Build(r, usecase)
	s := sub.CreateSub(app.cfg, *cl)

	sub, err := s.SubscribeToChannel("natschannel", stan.StartAtTime(time.Now()))

	defer sub.Unsubscribe()

	err = http.ListenAndServe("localhost:8181", r)
	if err != nil {
		log.Fatal("Сервер не запустился")
		return
	} else {
		log.Println("Server started")
	}
}
