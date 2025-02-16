# Avito-shop-backend

Реализация API для магазина мерча в рамках тестового задания для стажера Backend-направления.

## Основные возможности:
- Авторизация и получение JWT-токена
- Покупка товаров из внутреннего магазина
- Перевод монет между пользователями
- Просмотр баланса монет, инвентаря и истории транзакций


## Стек технологий
- Веб-фреймворк: **Gin**
- Аутентификация: **JWT**
- База данных: **PostgreSQL**
- ORM: **GORM**
- Деплой: **Docker**
- Тестирование: **Testify**
- Линтер: **golangci-lint**


## Установка и запуск

1. Клонировать репозиторий:
```bash
git clone https://github.com/PODTYAZHKI/avito_shop_backend && cd avito_shop_backend
```
2. Запустить сервисы:

```bash
docker-compose up
```  

3. Сервис доступен по адресу:
```
http://localhost:8080
```

## API

### 1. Аутентификация
**POST /api/auth**
- Регистрация или вход с именем пользователя и паролем
- ### request:
```json
{
  "username": "test",
  "password": "test"
}
```
- ### response:
```json
{
  "token": "jwt"
}
```

### 2. Перевод монет (protected)
**POST /api/sendCoin**
- Передача монет другому сотруднику
```
Authorization: Bearer <token>
```
- ### request:
```json
{
  "toUser": "anotherUser",
  "amount": 100
}
```
### 3. Покупка товара (protected)
**GET /api/buy/{item}**
```
Authorization: Bearer <token>
```

**Мерч** — это продукт, который можно купить за монетки. Всего в магазине доступно 10 видов мерча. Каждый товар имеет уникальное название и цену.

| Название     | Цена |
|--------------|------|
| t-shirt      | 80   |
| cup          | 20   |
| book         | 50   |
| pen          | 10   |
| powerbank    | 200  |
| hoody        | 300  |
| umbrella     | 200  |
| socks        | 10   |
| wallet       | 50   |
| pink-hoody   | 500  |

### 4. Получение информации (protected)
**GET /api/info**
- Получение баланса, истории транзакций и списка приобретенных товаров
```
Authorization: Bearer <token>
```
- ### response:
```json
{
    "coins": 850,
    "inventory": [
        {
            "type": "book",
            "quantity": 2
        }
    ],
    "coinHistory": {
        "received": null,
        "sent": [
            {
                "toUser": "testuser",
                "amount": 50
            }
        ]
    }
}
```
