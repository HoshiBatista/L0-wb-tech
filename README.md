# L0-wb-tech
# L0-wb-tech — Демонстрационный сервис заказов

[![CI Status](https://github.com/crissyro/L0-wb-tech/actions/workflows/go.yml/badge.svg)](https://github.com/crissyro/L0-wb-tech/actions/workflows/go.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Тестовый проект, разработанный в рамках обучения в WB Tech School. Это небольшой микросервис на Go, который получает данные о заказах из очереди сообщений, сохраняет их в базу данных и предоставляет доступ к ним через кэш и HTTP API.

---

### 💻 Стек технологий

<p>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Apache%20Kafka-231F20?style=for-the-badge&logo=apachekafka&logoColor=white" alt="Kafka">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/GitHub%20Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white" alt="GitHub Actions">
</p>

---

### ✨ Основные возможности

*   **Получение данных из Kafka:** Сервис подписывается на топик и получает сообщения с данными о заказах.
*   **Сохранение в PostgreSQL:** Валидные данные сохраняются в реляционную базу данных с использованием GORM.
*   **Кэширование в памяти:** Все заказы кэшируются в `map` для мгновенного доступа.
*   **Восстановление кэша:** При старте сервиса кэш автоматически заполняется данными из базы данных.
*   **HTTP API:** Реализован эндпоинт `GET /order/{order_uid}` для получения заказа из кэша.
*   **Веб-интерфейс:** Простая HTML-страница для поиска и просмотра информации о заказе по его ID.

---

### 🚀 Запуск и использование

#### Необходимые утилиты:
*   [Go](https://go.dev/) (версия 1.24+)
*   [Docker](https://www.docker.com/) и Docker Compose
*   [golangci-lint](https://golangci-lint.run/) (опционально, для локальной проверки)

#### Инструкция по запуску:

1.  **Клонируйте репозиторий:**
    ```bash
    git clone https://github.com/crissyro/L0-wb-tech.git
    cd L0-wb-tech
    ```

2.  **Запустите инфраструктуру (PostgreSQL и Kafka):**
    *Внимание: убедитесь, что Docker запущен на вашем компьютере.*
    ```bash
    docker-compose up -d
    ```

3.  **Запустите сервис на Go:**
    ```bash
    go run ./cmd/app/main.go
    ```
    После запуска вы увидите логи о подключении к БД, миграции и восстановлении кэша.

#### Как пользоваться:

*   **Веб-интерфейс:** Откройте в браузере [http://localhost:8081/](http://localhost:8081/).
*   **API:** Отправьте GET-запрос на `http://localhost:8081/order/{order_uid}`, чтобы получить данные в формате JSON.

---

###  linting

Для проверки качества кода и автоматического форматирования используется `golangci-lint`.

**Запустить проверку:**
```bash
golangci-lint run ./...
