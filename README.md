# Time-Tracker

### Инструкция по сборке:

1. Установите docker и docker-compose, если ранее
они не были установлены


2. Укажите url внешнего api. Для этого измените значение
переменной `EXTERNAL_API_URL` в файле `.env`


3. Включите docker и запустите контейнеры командой
    ```bash
    docker-compose up --build
    ```
    Сервер запустится на адресе
    ```bash
   http://localhost:8080/
   ```

### Документация swagger

Ознакомиться с документацией swagger можно
по ссылке:
```bash
http://localhost:8080/swagger/index.html
```

### P.S.
В базу данных добавлено 5 тестовых наборов данных
для пользователей
