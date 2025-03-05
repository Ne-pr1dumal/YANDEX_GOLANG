# YANDEX_GOLANG

Cервис для вычисления арифметических выражений через HTTP-запрос

## Структура проекта

```.
├── cmd
│   ├── agent_start
│   │   └── main.go
│   └── orchestrator_start
│       └── main.go
├── go.mod
└── internal
    ├── agent
    │   ├── agent.go
    │   └── agent_test.go
    └── orchestrator
        └── orchestrator.go
```
## Запуск проекта
0. Перейти в **главную** папку проекта (YANDEX_GOLANG)
1. Запуск orchestator`а

```bash
TIME_ADDITION_MS=200 TIME_SUBTRACTION_MS=200 TIME_MULTIPLICATIONS_MS=300 TIME_DIVISIONS_MS=400 go run cmd/orchestrator_start/main.go
```

**Ответ:**  Starting Orchestrator on port 8080.

3. Запуск Agent`а (в новом окне терминала)

```bash
COMPUTING_POWER=4 ORCHESTRATOR_URL=http://localhost:8080 go run cmd/agent_start/main.go
```

**Ответ:**
Starting Agent...
Starting worker 0
Starting worker 1
Starting worker 2
Starting worker 3

4. Отправка запроса (в новом окне терминала)
## Запросы:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '
{
  "expression": "2*2+2"
}'
```

**Ответ:** 
```{"id":"1"}```
Вашей заявке будет присвоен определенный **id**

Результ вычисленных операций предоставляется по данному запросу:

```bash
curl --location 'http://localhost:8080/api/v1/expressions'
```

Вывод:

```bash
{"expressions":[{"id":"1","expression":"2*2+2","status":"completed","result":6}]}
```

Можно выполнить поиск задачи по его id:

```bash
curl --location 'http://localhost:8080/api/v1/expressions/<id>'
```
где ```<id>``` - это номер заявки

# Примеры ошибок

**Ошибка на отсутствие выражения:**

```bash
{"error":"Expression not found"}
```

**Ошибка на невалидное выражение:**

```bash
{
  {"error":"expected number at position 1"}
}
```

**Ошибка на деление на ноль:**

```bash
{ Worker n: error computing task 3: division by zero }
```

## P. S.
```
Почему так мало коммитов?
Ответ прост: У меня сломался телефон с двухфакторкой 
Пришлось загрузить финалку на новом акке
```
# Credits
```
             _                   
  __ _  ___ | | __ _ _ __   __ _ 
 / _` |/ _ \| |/ _` | '_ \ / _` |
| (_| | (_) | | (_| | | | | (_| |
 \__, |\___/|_|\__,_|_| |_|\__, |
 |___/                     |___/
```
