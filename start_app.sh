#!/bin/bash

# Найти PID процесса с именем "main", использующего порт 8080
PID=$(sudo lsof -t -i :8080 -sTCP:LISTEN | xargs ps -o pid=,comm= | grep "main" | awk '{print $1}')

# Проверить, найден ли PID
if [ -n "$PID" ]; then
    echo "Завершение процесса 'main' с PID: $PID"
    sudo kill -9 $PID
else
    echo "Нет процесса 'main', использующего порт 8080."
fi

# Запуск Go-программы
echo "Запуск Go-программы..."