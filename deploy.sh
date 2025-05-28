#!/bin/bash

echo "📥 Обновление кода из git..."
git pull origin master || { echo "❌ git pull не удался"; exit 1; }

echo "⛔ Остановка и удаление старых контейнеров и volume..."
docker compose down -v

echo "🔨 Пересборка и запуск..."
docker compose up -d --build

echo "✅ Готово. Логи:"
docker compose logs -f