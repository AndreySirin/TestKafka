# Микросервис обработки сообщений

[Ссылка на тестовое задание](https://docs.google.com/document/u/0/d/13JHrzO9HuExWe_X0WrJzD8TPk3UHyTowUHCJHA5RFuY/mobilebasic?pli=1)

## Установка
```shell
git clone git@github.com:AndreySirin/TestKafka.git
```
## Запуск
```shell
docker compose up
```
## Удаление контейнеров
```shell
docker compose down
```

# Методы API
## Отправка сообщения
```shell
метод:POST
URL:http://localhost:8080/api/v1/message
Body:
{
    "message":"привет"
}
ответ: 
stutus 201
```

## Проверка статистики
```shell
метод:GET
URL:http://localhost:8080/api/v1/statistics
Body:
{}
ответ: 
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "__name__": "kafka_messages_sent_total",
          "instance": "app:8080",
          "job": "prometheus"
        },
        "value": [
          1751567928.454,
          "1"
          ]}]}}

stutus 200
```

