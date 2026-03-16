# Copilot Instructions – Stock Swipe App

## Project Overview

A high-performance web app for discovering S&P 500 stocks via swipe interaction.
Built on four decoupled layers: Python data worker → PostgreSQL → Go API → React UI.

**Core rule: No live external API calls during user interaction. All data is pre-fetched and stored.**

---

## Architecture

```
Python Worker  →  PostgreSQL  →  Go API  →  React UI
(data ingest)     (storage)      (serve)     (render)
```

---

## Layer 1 – Python Worker (`data-worker/`)

**Purpose:** Fetch, normalize, and store financial data. Never called by the API directly.

- Source: `yfinance` (Yahoo Finance)
- Runs on a schedule (cron / scheduler)
- All writes use UPSERT (never duplicate rows)
- All writes are transactional

**Key modules:**
- `tickers_loader.py` – load S&P 500 ticker list
- `fundamentals_fetcher.py` – PE, EPS, revenue growth, debt, sector
- `price_fetcher.py` – latest price + daily change %
- `history_fetcher.py` – OHLCV time-series
- `db_writer.py` – PostgreSQL UPSERT logic
- `scheduler.py` – job scheduling

**Refresh intervals:**
- Prices: every 5–15 minutes
- Fundamentals: once per day (nightly)
- Historical candles: daily append

---

## Layer 2 – PostgreSQL (Database)

All tables use UUIDs as primary keys. Use composite indexes on foreign key + date columns.

**Core tables:**

| Table | Purpose |
|---|---|
| `stocks` | Static + semi-static fundamentals |
| `stock_prices` | Latest price snapshot |
| `stock_history` | OHLCV time-series for charts |
| `users` | User accounts |
| `swipes` | Swipe history (left/right) per user |
| `watchlist` | Saved stocks per user |

**Key constraints:**
- `watchlist`: unique on `(user_id, stock_id)`
- `stock_history`: composite index on `(stock_id, date)`
- `swipes`: index on `user_id` and `(user_id, stock_id)`

---

## Layer 3 – Go API (`internal/`)

**Purpose:** Serve pre-computed data from PostgreSQL. No scraping. No external HTTP calls.

**Layered architecture — strict dependency direction:**

```
handler → service → repository → domain
```

| Layer | Responsibility |
|---|---|
| `domain/` | Pure Go structs/models. No DB, no HTTP. |
| `repository/` | SQL queries only. Uses `pgx` + `sqlc`. Returns domain models. |
| `service/` | Business logic: filtering, discover algorithm, watchlist rules. |
| `handler/` | HTTP routing, JSON serialization, input validation. |

**Core endpoints:**

```
GET  /stocks/discover              # Returns 1 stock (excludes swiped + watchlist)
POST /stocks/{id}/swipe            # Body: { "direction": "left" | "right" }
GET  /stocks/{ticker}              # Fundamentals + latest price
GET  /stocks/{ticker}/history?range=1y  # Time-series for chart
```

**Performance targets:** < 50ms from DB, < 150ms total per request.

**Rules:**
- Never call external APIs
- Never perform scraping
- Keep handlers thin — logic belongs in service layer
- Use `pgx` for DB connections, `sqlc` for generated queries

---

## Layer 4 – React Frontend (`features/`)

**Purpose:** UI only. No business logic on the client.

**Stack:**
- React + TypeScript
- Tailwind CSS
- Framer Motion (swipe gestures + animations)
- Zustand (state management)
- TradingView Lightweight Charts (price charts)

**Structure:**
```
features/
  swipe/
  watchlist/
  stockDetails/
shared/
  components/
  ui/
```

**Swipe flow:**
1. `GET /stocks/discover` → render card
2. User swipes → `POST /stocks/{id}/swipe`
3. Animate next card

No heavy computation on the client. All data comes from the API.

---

## Deployment

Docker Compose services:
- `postgres`
- `go-api`
- `python-worker`
- `frontend`

All secrets via environment variables (DB credentials, refresh intervals, API base URL). Never hardcode credentials.

---

## General Coding Guidelines

- **Separation of concerns:** Each layer does exactly one thing. Do not mix responsibilities.
- **Database-driven:** The API only reads what the worker already wrote. Never bypass this.
- **Deterministic responses:** API responses must be fast and predictable — no side effects on GET requests.
- **UPSERT everywhere in the worker:** Never assume a record doesn't exist.
- **Indexes matter:** Always add indexes for any column used in WHERE, JOIN, or ORDER BY.
- **Go:** Follow standard project layout (`internal/`, `cmd/`). Keep `main.go` minimal.
- **Python:** Keep each worker module focused on a single responsibility.
- **React:** Components are presentational. Data fetching lives in hooks or service files, not inside components.