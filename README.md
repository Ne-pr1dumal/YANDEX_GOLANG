# YANDEX_GOLANG или убийца PHOTOMATH

Cервис для вычисления арифметических выражений через HTTP-запрос

## Структура проекта

```.
.
├── README.md
├── cmd
│   ├── agent_start
│   │   ├── Dockerfile
│   │   └── main.go
│   └── orchestrator_start
│       ├── Dockerfile
│       └── main.go
├── docker-compose.yml
├── frontend
│   ├── Dockerfile
│   ├── nginx
│   │   └── nginx.conf
│   ├── package.json
│   ├── public
│   │   └── index.html
│   └── src
│       ├── App.js
│       ├── api.js
│       ├── index.css
│       └── index.js
├── go.mod
└── internal
    ├── agent
    │   ├── agent.go
    │   └── agent_test.go
    └── orchestrator
        └── orchestrator.go
```
## Скачайте проект

1. Склонируйте проект с GitHub
    ```bash
    git clone https://github.com/Ne-pr1dumal/Calculation-Service-Yandex
    ```
2. Перейдите в головную папку с проектом и запустите проект
    ```bash
    go run ./cmd/main.go
    ```
## Запуск проекта

0. Перейти в **главную** папку проекта (YANDEX_GOLANG)
1. Вводим запрос для запуска сервисов
```bash
docker-compose up
```

В терминале увидим:
```
agent-1         | 2025/03/06 15:50:43 Agent is Starting...
agent-1         | 2025/03/06 15:50:43 Starting worker 0
agent-1         | 2025/03/06 15:50:43 Starting worker 1
agent-1         | 2025/03/06 15:50:43 Starting worker 2
agent-1         | 2025/03/06 15:50:43 Starting worker 3
```
2. В новом терминале делаем http-запрос
(Например)
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '
{
  "expression": "2*2+2"
}'
```
4. В окне с Docker`ом увидим
```
agent-1         | 2025/03/06 08:41:42 Worker 3: processing task 1: 2.000000 * 2.000000 (300 ms)
agent-1         | 2025/03/06 08:41:42 Worker 3: successfully completed task 1 with result 4.000000
agent-1         | 2025/03/06 08:41:42 Worker 3: processing task 2: 4.000000 + 2.000000 (200 ms)
agent-1         | 2025/03/06 08:41:43 Worker 3: successfully completed task 2 with result 6.000000
```
5. Также есть фронтенд (Перейдите по ссылке)
```
http://localhost:3000
```
6. В поле ввести выражение
7. Результат можно увидеть только при отправке нового выражения (F5 или CMD+R работать не будут)
8. Кайфуем

![Image](https://github.com/Ne-pr1dumal/YANDEX_GOLANG/blob/experemental/Снимок%20экрана%202025-03-06%20в%2018.58.11.png)

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

## Примеры ошибок

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
## Также есть запуск без Docker
### Обычный
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
5. Кайфуем

# FAQ

**Оркестратор**:

- Принимает выражения через REST API
- Разбивает выражения на атомарные задачи
- Управляет очередью задач
- Собирает результаты
- Хранит статусы вычислений

**Агенты**:

- Получают задачи через HTTP-запросы
- Выполняют арифметические операции с задержкой
- Возвращают результаты через API

# P. S.
```
Почему так мало коммитов?
Ответ прост: У меня сломался телефон с двухфакторкой 
Пришлось загрузить финалку на новом акке
```
# Credits
```
 ____  _   _  ___ _____ ___  __  __    _  _____ _   _ 
|  _ \| | | |/ _ \_   _/ _ \|  \/  |  / \|_   _| | | |
| |_) | |_| | | | || || | | | |\/| | / _ \ | | | |_| |
|  __/|  _  | |_| || || |_| | |  | |/ ___ \| | |  _  |
|_|   |_| |_|\___/ |_| \___/|_|  |_/_/   \_\_| |_| |_|
                                                      
 __  __ _   _ ____  ____  _____ ____  _____ ____  
|  \/  | | | |  _ \|  _ \| ____|  _ \| ____|  _ \ 
| |\/| | | | | |_) | | | |  _| | |_) |  _| | |_) |
| |  | | |_| |  _ <| |_| | |___|  _ <| |___|  _ < 
|_|  |_|\___/|_| \_\____/|_____|_| \_\_____|_| \_\
                                      _   _           
 _ __   _____      _____ _ __ ___  __| | | |__  _   _ 
| '_ \ / _ \ \ /\ / / _ \ '__/ _ \/ _` | | '_ \| | | |
| |_) | (_) \ V  V /  __/ | |  __/ (_| | | |_) | |_| |
| .__/ \___/ \_/\_/ \___|_|  \___|\__,_| |_.__/ \__, |
|_|                                             |___/ 
             _                   
  __ _  ___ | | __ _ _ __   __ _ 
 / _` |/ _ \| |/ _` | '_ \ / _` |
| (_| | (_) | | (_| | | | | (_| |
 \__, |\___/|_|\__,_|_| |_|\__, |
 |___/                     |___/
```
