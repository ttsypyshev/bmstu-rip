# RIP (разработка интернет приложений)

Приложение выполнено на `go 1.23.2`
[Ссылка на курс](https://github.com/iu5git/Web)

## Параметры запуска

Запустите приложение с помощью следующих команд:

```sh
sudo docker-compose up -d
chmod +x start_app.sh
./start_app.sh && go run src/main.go
```

## Настройка базы данных

Для работы с базой данных используйте следующие SQL-скрипты:

- [Создать базу данных](database/create.sql)
- [Вставить данные в базу данных](database/insert.sql)
- [Удалить базу данных](database/drop.sql)