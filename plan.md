# План: Go‑бот Freyja (FSD), полный паритет с n8n

## Краткое резюме
Собрать Go‑сервис (Go `1.25.7`) по FSD‑подходу внутри `internal/`, использующий `telebot` (long polling), `sqlx` (Supabase Postgres), OpenRouter (оба LLM‑сценария), Supabase Storage (S3‑совместимо), с логированием сообщений бота в историю и фоновым удалением по расписанию.

## Важные изменения/публичные интерфейсы
1. Публичные команды Telegram: `/find_card_scan`, `/draw_card`.
2. Публичные конфиги через ENV (см. ниже).
3. Поведение в группах: реакции только на сообщения с упоминанием бота.
4. Отдельный фоновый процесс удаления сообщений каждые `1m`.

## Версия Go
1. `go.mod`:
   1. `go 1.25.7`
   2. `toolchain go1.25.7`

## Явный список библиотек
1. Telegram: `gopkg.in/telebot.v3`
2. Postgres: `github.com/jmoiron/sqlx`
3. PostgreSQL драйвер: `github.com/jackc/pgx/v5/stdlib`
4. Логи: `github.com/sirupsen/logrus`
5. Env: `github.com/joho/godotenv`
6. S3 (Supabase Storage): `github.com/aws/aws-sdk-go-v2`
7. HTTP/JSON: стандартные `net/http`, `encoding/json`

## FSD‑структура проекта (внутри `internal/`)
1. `cmd/bot/main.go`
2. `internal/app/`
   1. `bootstrap.go` — создание конфига, логгера, клиентов, зависимостей.
   2. `router.go` — регистрация telebot‑хэндлеров.
3. `internal/shared/`
   1. `config/` — `godotenv` + структура конфига.
   2. `logger/` — настройка `logrus`.
   3. `telegram/` — настройка `telebot`, общие helpers (send/delete/log).
   4. `db/` — подключение `sqlx.DB`.
   5. `storage/` — S3‑клиент для Supabase Storage.
   6. `llm/` — OpenRouter клиент + парсер JSON.
   7. `time/` — helpers для `delete_at`.
4. `internal/entities/`
   1. `card/` — модели + `repo` + `repo_sqlx`.
   2. `drawn/` — модели + `repo` + `repo_sqlx`.
   3. `history/` — модели + `repo` + `repo_sqlx`.
5. `internal/features/`
   1. `find_card_scan/` — `handler.go`, `service.go`, `llm_prompt.go`.
   2. `draw_card/` — `handler.go`, `service.go`, `llm_prompt.go`.
6. `internal/processes/`
   1. `delete_messages/runner.go`.
7. `go.mod` с `module neurofreyja`.

## Режим Telegram
1. Только long polling.
2. `telebot.Settings` с `LongPoller{Timeout: 10s}`.
3. `HTTPClient` с таймаутом `30s`.

## ENV‑конфигурация
1. `TELEGRAM_TOKEN`
2. `BOT_USERNAME` (если пусто — берём через `GetMe`).
3. `OPENROUTER_API_KEY`
4. `OPENROUTER_BASE_URL` (default `https://openrouter.ai/api/v1`)
5. `OPENROUTER_MODEL_TITLES` (default `openai/gpt-4.1`)
6. `OPENROUTER_MODEL_DESC` (default `openai/gpt-4.1`)
7. `PG_DSN`
8. `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`
9. `S3_BUCKET` (default `neurofreyja`)
10. `S3_REGION`
11. `S3_PATH_STYLE` (default `true`)
12. `DELETE_AFTER_MINUTES` (default `10`)
13. `DELETE_SWEEP_INTERVAL` (default `1m`)

## База данных (существующая)
1. `neuro_freyja_card(id, title, image_id, description, created_at, updated_at)`
2. `neuro_freyja_drawn_card(id, card_id, chat_id, created_at)`
3. `neuro_freyja_history(id, message_id, chat_id, chat_title, chat_type, sent_at, content, delete_at, deleted_at)`

## Командный роутинг (паритет n8n)
1. `message_text` берётся из `text` или `caption`.
2. Тип чата:
   1. `group`/`supergroup`: реагируем только если есть `@BOT_USERNAME`.
   2. `private`: всегда реагируем.
3. Удаляем `@BOT_USERNAME` из текста и `TrimSpace`.
4. Роутинг:
   1. `"/find_card_scan"` → фича `find_card_scan`.
   2. `"/draw_card"` → фича `draw_card`.
   3. Иначе → ответ `Такая команда мне неизвестна` + лог в историю.

## Фича `/find_card_scan`
1. Отправить `Сейчас поищу...`.
2. Проверка reply:
   1. Нет reply → `Эту команду надо вызывать реплаем на сообщение` + лог, завершить.
3. Взять текст из reply (`caption` приоритетнее).
4. LLM (OpenRouter, `OPENROUTER_MODEL_TITLES`) с промптом из n8n, ожидается JSON `{"titles":[...]}`.
5. Парсинг: строгий JSON, при мусоре — извлечь первый JSON‑объект и парсить.
6. Нормализация: `strings.ToLower`, `TrimSpace`, убрать пустые, `dedupe`.
7. SQL: `select title, image_id from neuro_freyja_card where lower(title) = any($1)` (`text[]`).
8. Для каждой найденной карты:
   1. Скачать изображение из S3 по `image_id` (context timeout `30s`).
   2. Отправить фото с подписью `Это карта - "title"`.
   3. Отправить предупреждение `У вас есть 10 минут чтобы скачать изображения`.
   4. Логировать оба сообщения с `delete_at = now()+10m`.
9. Удалить сообщение `Сейчас поищу...` в конце (даже при ошибках).

## Фича `/draw_card`
1. Отправить `Вытягиваю карту...`.
2. SQL:
   1. `select c.id, c.title, c.image_id, c.description from neuro_freyja_card c where c.id not in (select card_id from neuro_freyja_drawn_card where chat_id=$1) order by random() limit 1`.
3. Нет карты → удалить загрузочное и завершить.
4. Если `description` пустое:
   1. LLM (OpenRouter, `OPENROUTER_MODEL_DESC`) с промптом из n8n.
   2. `update neuro_freyja_card set description=$1, updated_at=now() where id=$2`.
5. Скачать изображение из S3 по `image_id`.
6. Отправить фото с подписью `Карта "<b>{title}</b>"\n\n{description}`, `parse_mode=HTML`.
7. Логировать сообщение (без `delete_at`).
8. `insert into neuro_freyja_drawn_card(card_id, chat_id) values ($1,$2)`.
9. Удалить сообщение `Вытягиваю карту...` в конце.

## История сообщений
1. Сохраняем все исходящие сообщения бота.
2. `sent_at = time.Unix(msg.Date,0).UTC()`.
3. `content = msg.Caption` иначе `msg.Text`.

## Фоновое удаление сообщений
1. Тикер каждые `DELETE_SWEEP_INTERVAL`.
2. SQL:
   1. `select id, message_id, chat_id from neuro_freyja_history where delete_at is not null and delete_at <= now() and deleted_at is null`.
3. Удаляем через Telegram.
4. `update neuro_freyja_history set deleted_at=now() where id = any($1)`.

## LLM и таймауты
1. OpenRouter запросы через `net/http` с `Timeout: 60s`.
2. DB запросы с контекстом `5s`.
3. S3 загрузка с контекстом `30s`.

## Тесты и сценарии
1. Роутинг:
   1. group/supergroup с и без упоминания.
   2. private без упоминания.
2. `/find_card_scan`:
   1. Нет reply → корректный ответ.
   2. Reply с caption/text → корректный LLM запрос и SQL.
   3. Нет результатов → нет дополнительных сообщений, удаляется только загрузочное.
   4. Несколько карт → фото + предупреждения, оба логируются с `delete_at`.
3. `/draw_card`:
   1. Есть `description` → без LLM.
   2. Нет `description` → LLM + update.
   3. Нет карт → только удаление загрузочного.
4. Sweep:
   1. Удаление сообщений по `delete_at`.
   2. Ошибка удаления → запись остаётся для повторной попытки.

## Допущения и дефолты
1. `.env` может отсутствовать — запуск продолжается.
2. Поведение ошибок: только лог, без дополнительных пользовательских сообщений.
3. Предупреждение об удалении отправляется для каждого изображения.
4. Миграции не добавляются, используются существующие таблицы Supabase.
