# 💳 Tochka Medusa Integration

Интеграция платёжной платформы **Точка (Medusa API)** для эскроу-флоу: вознаграждение + депозит в рамках сделок между пользователями.

---

## 📌 Общая схема

1. Обе стороны подтверждают участие в сделке
2. Создаются два платёжных заказа в Точке:
   - Заказчик платит вознаграждение
   - Перевозчик вносит страховой депозит
3. Точка удерживает средства (эскроу)
4. После завершения сделки:
   - Деньги выплачиваются или возвращаются, в зависимости от статуса

---

## ✅ Подключение к Точке

### 1. Зарегистрироваться
Обратитесь в поддержку Точка Банка:
- Укажите ИНН, юр.лицо/ИП, email
- Уточните, что хотите подключить **Medusa API**

### 2. Сгенерировать ключи

```bash
openssl genrsa -out rsaprivkey.pem 2048
openssl req -days 1825 -new -x509 -key rsaprivkey.pem -out rsaacert.pem
```

Передайте банку:
- `rsaacert.pem` (публичный ключ)
- Акт и соглашение (шаблоны даст менеджер)

### 3. Получите от Точки:
- `username`, `password` (для Basic Auth)
- `Sign-Key-Id`
- URL: `https://uapi.tochka.com/uapi/medusa/v1.0`

---

## 🔧 Используемые эндпоинты

| Категория | Описание | Метод |
|----------|----------|-------|
| Получатель | Создать получателя | `POST /recipients` |
| Карта | Добавить карту | `POST /recipients/{id}/payout_methods/cards` |
| Заказ | Создать заказ | `POST /orders` |
| Заказ | Проверить статус | `GET /orders/{orderExtId}` |
| Решение | Подтвердить/отклонить услугу | `POST /orders/{orderExtId}/decisions` |
| Выплата | Выплатить получателю | `POST /proceed_service_payout_to_recipient` |
| Возврат | Вернуть плательщику | `POST /proceed_refund` |

---

## 🔁 Жизненный цикл сделки

### 1. Сделка подтверждена
- Создаём 2 заказа (`POST /orders`)
  - `orderExtId = deal_42_reward`
  - `orderExtId = deal_42_deposit`

### 2. Получаем `paymentUrl`, отправляем юзерам

### 3. Проверяем оплату через `GET /orders/{orderExtId}`
- `"state": "PAID_BY_ACQUIRER"` → отметка `isPaid = true`

### 4. После выполнения:
- `POST /orders/{id}/decisions` → `"confirmed"`
- `POST /proceed_service_payout_to_recipient`

### 5. Если провал:
- `POST /orders/{id}/decisions` → `"rejected"`
- `POST /proceed_refund`

---

## 🔒 Безопасность

- Используйте RSA-подпись для всех запросов, где требуется
- Подпись уходит в `Sign-Body` (base64)
- `Sign-Key-Id` — ID, выданный Точкой

---

## 🗃️ Что хранить в БД

- `recipientId` для каждого пользователя
- `cardExtId` (если используется)
- `orderExtId` (на каждый платёж)
- `paymentUrl`
- текущий статус оплаты
- результат выплаты/возврата
- тип: `reward` / `deposit`

---

## 📎 Пример orderExtId

```
deal_42_reward
deal_42_deposit
```

---

## 👀 Примечания

- Один заказ = один платёж = один получатель
- Нельзя сделать выплату "третьей стороне" без указания её в `Recipient`
- Вебхуки необязательны — можно опрашивать вручную

---

## 🆘 Тестирование

Точка предоставляет **песочницу** с имитацией всех этапов:
- успешная оплата
- ошибка эквайринга
- возврат
- комиссия

Песочные методы — через `/sandbox/...`, аналогичные боевым.
