# Stellar Mass Calculator - REST API

REST API веб-приложения для расчета масс звезд по спектральным классам

## Домены API

### Домен спектрального класса (SpectralClass)

**GET /api/classes** - Список спектральных классов с фильтрацией  
**GET /api/classes/:id** - Получение спектрального класса по ID  
**POST /api/classes** - Создание нового спектрального класса  
**PUT /api/classes/:id** - Обновление спектрального класса  
**DELETE /api/classes/:id** - Удаление спектрального класса

### Домен расчета масс (MassRequest)

**GET /api/mass-requests/star-calculation** - Информация о текущей заявке  
**GET /api/mass-requests** - Список заявок с фильтрацией  
**GET /api/mass-requests/:id** - Получение заявки по id
**POST /api/mass-requests** - Создание новой заявки  
**PUT /api/mass-requests/:id** - Обновление полей заявки  
**PUT /api/mass-requests/:id/form** - Формирование пустой заявки 
**PUT /api/mass-requests/:id/complete** - Завершение заявки
**DELETE /api/mass-requests/:id** - Удаление заявки

### Домен связь расчет-класс (MassRequest-Class)

**POST /api/mass-requests/classes** - Добавление спектрального класса в заявку  
**PUT /api/mass-requests/:id/classes/:class_id** - Изменение параметров класса в заявке  
**DELETE /api/mass-requests/:id/classes/:class_id** - Удаление класса из заявки

### Домен пользователь (User)

**POST /api/auth/register** - Регистрация пользователя  
**POST /api/auth/login** - Аутентификация  
**POST /api/auth/logout** - Деавторизация  
**GET /api/user/profile** - Получение профиля пользователя  
**PUT /api/user/profile** - Обновление профиля пользователя

## Модели данных

### Спектральный класс (Class)

```json
{
    "ID": 1,
    "Name": "A",
    "Temperature": 7500,
    "Color": "Белый",
    "Spectre": "A0-A9",
    "Examples": "Сириус, Вега",
    "IsDeleted": false,
    "Image": "/images/A-Class.png",
    "MassRequestToClass": null
}
```

```json
{
    "ID": 0,
    "Status": 1,
    "UserID": 1,
    "User": {
        "ID": 1,
        "Username": "newusername",
        "PassHash": "newpassword",
        "PassSalt": "",
        "IsMod": false
    },
    "ModeratorId": 1,
    "Moderator": {
        "ID": 1,
        "Username": "newusername",
        "PassHash": "newpassword",
        "PassSalt": "",
        "IsMod": false
    },
    "MassRequestToClass": null,
    "CreatedAt": "2025-10-15T03:10:59.925301+03:00",
    "FormedAt": "2025-10-15T03:10:59.924016+03:00",
    "ClosedAt": "2025-10-15T03:10:59.924016+03:00"
}
```
## Статусы расчетов
1 - Черновик (Draft)
2 - Удален (Deleted)
3 - Ожидание (Pending)
4 - Завершено (Completed)
5 - Отклонено (Rejected)

## Бизнес-логика

### Расчет массы звезды

При завершении расчета вычисляется масса звезды по формуле:
**M = L^exponent**

**Где:**
- M - масса звезды (в солнечных массах)
- L - светимость звезды (в солнечных светимостях)
- exponent - показатель степени в зависимости от спектрального класса

### Показатели степени для разных классов:
O, B, A классы: M = L^0.222
F, G классы: M = L^0.234
K, M классы: M = L^0.264
Прочие классы: M = L^0.25


### Особенности расчета

- Расчетные данные (массы звезд) сохраняются в базе данных
- Светимость измеряется в единицах солнечной светимости
- Масса возвращается в солнечных массах с точностью до 0.01
- Пользователь по умолчанию имеет ID = 1

## Workflow статусов
Астроном: создает черновик → добавляет спектральные классы → формирует расчет
Модератор: отклоняет или завершает сформированный расчет с вычислением масс

- Пустой расчет-черновик создается автоматически при добавлении первого класса
- Системные поля вычисляются автоматически на бэкенде