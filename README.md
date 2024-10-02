# MusicLibrary

MusicLibrary — это веб-сервис для управления музыкальной библиотекой.

## Установка

1. Клонируйте репозиторий
    ```
    git clone https://github.com/Ivzond/MusicLibrary.git
    cd MusicLibrary
    ```
2. Установите зависимости
    ```
    go mod tidy 
    ```
3. Создайте файл .env в корне проекта со следующими переменными
    ```
    DATABASE_URL=<ваш_путь_к_базе_данных>
    PORT=<порт_для_сервера>
    EXTERNAL_API_URL=<url_внешнего_API>
    ```

## Запуск
1. Для запуска сервера выполните:
    ```
    go run cmd/server/main.go
    ```
   
## Примечание:

*Файл .env я добавил в .gitignore, но по умолчанию его можно заполнить, например, так:*
```
DATABASE_URL=postgres://postgres:12345678@localhost:5432/music-library?sslmode=disable
EXTERNAL_API_URL=track_info.com
```

*Так или иначе нужно будет изменить url БД и путь к внешнему API* 

*Также можно явно прописать номер порта, на котором поднимется сервер(по умолчанию :8080)*

