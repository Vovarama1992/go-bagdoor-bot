# Переменные
ENV_FILE=.env
COMPOSE=docker-compose --env-file $(ENV_FILE)

# Сервисы
MIGRATOR=migrator
BOT=bot
SERVER=server
DB=db

# Команды

up: ## Запустить все контейнеры (без запуска миграций)
	$(COMPOSE) up -d $(DB) $(BOT) $(SERVER)

migrate: ## Запустить миграции
	$(COMPOSE) run --rm $(MIGRATOR)

up-dev: migrate up ## Запустить миграции и потом все сервисы

stop: ## Остановить все контейнеры
	$(COMPOSE) down

logs: ## Смотреть логи всех контейнеров
	$(COMPOSE) logs -f

logs-bot: ## Смотреть логи только бота
	$(COMPOSE) logs -f $(BOT)

rebuild: ## Пересобрать все контейнеры
	$(COMPOSE) build --no-cache

ps: ## Показать состояние контейнеров
	$(COMPOSE) ps

restart-bot: ## Перезапустить только bot
	$(COMPOSE) restart $(BOT)

restart-server: ## Перезапустить только server
	$(COMPOSE) restart $(SERVER)

.PHONY: up migrate up-dev stop logs logs-bot rebuild ps restart-bot restart-server