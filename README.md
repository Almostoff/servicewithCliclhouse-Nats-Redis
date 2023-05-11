## Развернуть сервис на Golang, Postgres, Clickhouse, Nats (альтернатива kafka), Redis

### Описать модели данных и миграций

##### В миграциях Postgres

- Проставить primary-key и индексы на указанные поля

- При добавлении записи в таблицу устанавливать приоритет как макс приоритет в таблице +1. Приоритеты начинаются с 1

- При накатке миграций добавить одну запись в Campaigns таблицу по умолчанию


##### Реализовать CRUD методы на GET-POST-PATCH-DELETE данных в таблице ITEMS в Postgres

- При редактировании данных в Postgres ставить блокировку на чтение записи и оборачивать все в транзакцию. Валидируем поля при редактировании.

- При редактировании данных в ITEMS инвалидируем данные в REDIS

- Если записи нет (проверяем на PATCH-DELETE), выдаем ошибку (статус 404)

- При GET запросе данных из Postgres кешировать данные в Redis на минуту. Пытаемся получить данные сперва из Redis, если их нет, идем в БД и кладем их в REDIS

- При добавлении, редактировании или удалении записи в Postgres писать лог в Clickhouse через очередь Nats (альтернатива kafka). Логи писать пачками в Clickhouse


## Запуск

- Собрать докер-контейнер в родительском каталоге командой: docker-compose up
- Запустить server.go в каталоге server командой go run server.go

_ для удобства все тесты прописаны в файле test.http

__ P.s. развернут на локальном хосте//
