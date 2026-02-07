# Freyja Bot (Go)

Go service that mirrors the n8n workflow for `/find_card_scan` and `/draw_card`.

## Requirements
- Go `1.25.7`
- Postgres (Supabase)
- S3-compatible storage (Supabase Storage)
- Telegram bot token
- OpenRouter API key

## Setup
1. Create `.env` based on `.env.example`.
2. Run the bot:

```bash
cd /Users/ilya/projects/neurofreyja
# Optional: load .env if your shell doesn't
# set -a && source .env && set +a

go run ./cmd/bot
```

## Commands
- `/find_card_scan` (reply to a message)
- `/draw_card`

## Notes
- In group chats, the bot responds only when mentioned.
- Outgoing messages are logged to `neuro_freyja_history` and deleted on schedule if `delete_at` is set.
