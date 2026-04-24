# bridge-taiga-matrix

Лёгкий Go-сервис, который отслеживает дедлайны в проектах [Taiga](https://taiga.io) и отправляет уведомления в комнаты [Matrix](https://matrix.org).

## Принцип работы

При каждом запуске сервис:
1. Аутентифицируется в Taiga API по логину и паролю.
2. Перебирает настроенные проекты.
3. Получает все открытые пользовательские истории и задачи по каждому проекту.
4. Для каждого элемента проверяет:
   - дедлайн не задан → отправляет предупреждение;
   - дедлайн уже прошёл → отправляет уведомление о просрочке;
   - до дедлайна осталось не более `daysUntilDeadline` дней → отправляет напоминание.
5. Отправляет сообщение в Matrix-комнату проекта и, при необходимости, в общую комнату.

Предназначен для запуска по расписанию (cron, Планировщик задач Windows, systemd timer и т.п.).

## Требования

- Go 1.22+
- Доступ к экземпляру Taiga (self-hosted или облачный)
- Учётная запись Matrix с действующим токеном доступа и членством в нужных комнатах

## Настройка

Скопируйте `settings-example.json` в `settings.json` и заполните значения.

```json
{
    "taigaBaseURL"  : "https://taiga.example.com",
    "taigaUsername" : "ваш-логин",
    "taigaPassword" : "ваш-пароль",
    "taigaProjects" : [
        {
            "name": "slug-проекта",
            "matrixProjectRoomID": "!roomid:matrix.example.com"
        }
    ],

    "matrixServer" : "https://matrix.example.com",
    "matrixToken"  : "syt_ваш_matrix_токен",
    "duplicateToGeneralGroup": true,
    "generalRoomId": "!generalroomid:matrix.example.com",

    "InsecureSkipVerify": false,

    "language": "ru",

    "daysUntilDeadline": 15
}
```

### Параметры

| Параметр | Тип | Описание |
|---|---|---|
| `taigaBaseURL` | string | Базовый URL вашего экземпляра Taiga |
| `taigaUsername` | string | Логин в Taiga |
| `taigaPassword` | string | Пароль в Taiga |
| `taigaProjects` | array | Список проектов для мониторинга (см. ниже) |
| `taigaProjects[].name` | string | Slug проекта (виден в URL Taiga) |
| `taigaProjects[].matrixProjectRoomID` | string | ID Matrix-комнаты для этого проекта |
| `matrixServer` | string | Базовый URL вашего Matrix-сервера |
| `matrixToken` | string | Токен доступа Matrix |
| `duplicateToGeneralGroup` | bool | Дублировать все уведомления в `generalRoomId` |
| `generalRoomId` | string | ID общей Matrix-комнаты |
| `InsecureSkipVerify` | bool | Отключить проверку TLS-сертификата (в продакшене используйте `false`) |
| `language` | string | Язык уведомлений: `en` или `ru` (пустая строка — системная локаль) |
| `daysUntilDeadline` | int | Отправлять напоминание, если до дедлайна осталось столько дней или меньше |

## Сборка и запуск

```bash
go build -o bridge-taiga-matrix .
./bridge-taiga-matrix -config settings.json
```

Флаг `-config` необязателен; по умолчанию используется `settings.json` в текущей директории.

## Локализация

Шаблоны уведомлений хранятся в `locales/<lang>.json`. Чтобы добавить новый язык, создайте соответствующий файл и укажите его код в параметре `"language"` конфига.

## Пример планирования (Linux cron)

Запуск каждый день в 09:00:

```cron
0 9 * * * /opt/bridge-taiga-matrix/bridge-taiga-matrix -config /opt/bridge-taiga-matrix/settings.json
```
