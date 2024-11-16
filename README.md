# RIP (разработка интернет приложений)

Приложение выполнено на `go 1.23.2`
[Ссылка на курс](https://github.com/iu5git/Web)

## Параметры запуска

Запустите приложение с помощью следующих команд:

```sh
sudo docker compose up -d
chmod +x start_app.sh
./start_app.sh && go run src/main.go
```

## Ссылки для работы

1. **Панель просмотра хранилища**  
   Для просмотра содержимого бакета `code-inspector` используйте следующую ссылку:  
   [http://localhost:9001/buckets/code-inspector/browse](http://localhost:9001/buckets/code-inspector/browse)  
   Эта ссылка перенаправляет вас с (нужно создать бакет): `http://localhost:9000`.

2. **Браузер базы данных**  
   Для доступа к инструменту просмотра базы данных перейдите по ссылке (нужно создать сервер и таблицы):  
   [http://localhost:15432/browser/](http://localhost:15432/browser/)  


## Настройка базы данных

Для работы с базой данных используйте следующие SQL-скрипты:

- [Создать базу данных](database/create.sql)
- [Вставить данные в базу данных](database/insert.sql)
- [Удалить базу данных](database/drop.sql)