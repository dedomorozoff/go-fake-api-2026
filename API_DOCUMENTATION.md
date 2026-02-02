# Документация API для приложения интерактивных досок

## Аутентификация и регистрация

Все защищенные запросы должны содержать заголовок `Authorization: Bearer <token>`.

### Регистрация
`POST /registration`

**Запрос:**
```json
{
  "name": "Ivan",
  "email": "ivan@example.com",
  "password": "Password123!"
}
```
*Валидация:*
- `name`: только латиница.
- `password`: от 8 символов, должен содержать цифры и спецсимволы.

### Авторизация
`POST /authorization`

**Запрос:**
```json
{
  "email": "ivan@example.com",
  "password": "Password123!"
}
```

**Ответ:**
```json
{
  "data": {
    "user": {
      "id": 1,
      "name": "Ivan",
      "email": "ivan@example.com"
    },
    "token": "eyJhbGciOiJIUzI1Ni..."
  }
}
```

---

## Управление досками

### Создание доски
`POST /boards` (защищенный)

**Запрос:**
```json
{
  "name": "Моя новая доска",
  "is_public": true
}
```

---

### Список моих досок
`GET /boards` (защищенный)
Возвращает список досок, созданных пользователем или к которым ему предоставлен доступ.

---

### Предоставление доступа
`POST /boards/{board_id}/share` (защищенный)

**Запрос:**
```json
{
  "email": "friend@example.com"
}
```

---

### Лайк доске
`POST /boards/{board_id}/like` (защищенный)
Ставит или убирает лайк доске.

---

### Список публичных досок
`GET /public-boards`
Возвращает список досок с `is_public: true`, отсортированный по количеству лайков.

---

### Публичный просмотр доски
`GET /board/{hash}`
Доступ к доске по публичной ссылке (без авторизации).

---

## Работа в реальном времени (WebSocket)

Подключение: `ws://localhost:8080/ws/board/{board_id}?token=<token>`

### Формат сообщения (JSON)
```json
{
  "type": "string",
  "payload": "any"
}
```

### Типы сообщений (Client -> Server)

1. **Обновление объекта** (`object_update`):
   ```json
   {
     "type": "object_update",
     "payload": {
       "id": "obj1",
       "type": "rectangle",
       "x": 100,
       "y": 150,
       "width": 200,
       "height": 100,
       "rotation": 0,
       "color": "#ff0000"
     }
   }
   ```

2. **Захват фокуса** (`object_focus`):
   ```json
   {
     "type": "object_focus",
     "payload": "obj1"
   }
   ```
   *После этого объект блокируется для других пользователей.*

3. **Снятие фокуса** (`object_blur`):
   ```json
   {
     "type": "object_blur",
     "payload": "obj1"
   }
   ```

4. **Удаление объекта** (`object_delete`):
   ```json
   {
     "type": "object_delete",
     "payload": "obj1"
   }
   ```

### Сообщения от сервера (Server -> Client)
Сервер рассылает всем подключенным к доске те же сообщения, что получает от клиентов, добавляя информацию о пользователе (например, `owner_name` при захвате фокуса).
