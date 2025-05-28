#!/bin/bash

echo "⛔ Остановка и удаление старых контейнеров и volume..."
docker compose down -v

echo "🔨 Пересборка и запуск..."
docker compose up -d --build

echo "✅ Готово. Логи:"
docker compose logs -f