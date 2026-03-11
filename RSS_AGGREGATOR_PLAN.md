# RSS Aggregator API: Architectural Blueprint

## 1. System Components

### Fetcher (The Producer)

- Stateless service that reads `.yaml` files from a `/sources` directory
- Periodically polls RSS URLs
- Publishes a standard JSON payload to a RabbitMQ exchange for every feed item found
- **No database access**

### RabbitMQ (The Message Broker)

Acts as the buffer between data ingestion and data persistence.

| Resource           | Name              |
| ------------------ | ----------------- |
| Exchange           | `rss_events`      |
| Queue              | `raw_posts_queue`  |

### Processor (The Consumer/Worker)

- Listens to `raw_posts_queue`
- Performs de-duplication: checks if `guid` or `link` already exists in the DB
- **Primary Writer** — the only service with `INSERT`/`UPDATE` permissions on the database
- Syncs `/sources` YAML metadata into the `sources` table on startup

### API (The Server)

- Provides REST or GraphQL endpoints
- **Read-only** — only performs `SELECT` queries on the database
- Serves post data and the list of available sources

---

## 2. Source Configuration Format

Contributors add new feeds by creating a file in `/sources/`:

```yaml
# /sources/the-verge-tech.yaml
id: "the-verge-tech"
name: "The Verge"
url: "https://www.theverge.com/rss/index.xml"
category: "technology"
```
---

## 3. Database Schema (PostgreSQL)

### `sources`

| Column          | Type        | Constraints          |
| --------------- | ----------- | -------------------- |
| `id`            | `text`      | PK (slug from YAML)  |
| `name`          | `text`      | NOT NULL              |
| `url`           | `text`      | NOT NULL              |
| `category`      | `text`      |                      |
| `last_polled_at`| `timestamptz` |                    |

### `posts`

| Column          | Type        | Constraints              |
| --------------- | ----------- | ------------------------ |
| `id`            | `uuid`      | PK, default `gen_random_uuid()` |
| `source_id`     | `text`      | FK &rarr; `sources.id`  |
| `guid`          | `text`      | UNIQUE INDEX             |
| `title`         | `text`      | NOT NULL                 |
| `link`          | `text`      | NOT NULL                 |
| `description`   | `text`      |                          |
| `published_at`  | `timestamptz` |                        |

---

## 4. API Interface

| Method | Endpoint                          | Description                        |
| ------ | --------------------------------- | ---------------------------------- |
| GET    | `/sources`                        | List all sources from the DB       |
| GET    | `/posts`                          | Aggregated posts (paginated)       |
| GET    | `/posts?source_id=the-verge-tech` | Posts filtered by source            |

---

## 5. RabbitMQ Message Schema

Canonical JSON payload published by the Fetcher and consumed by the Processor:

```json
{
  "source_id": "the-verge-tech",
  "guid": "https://www.theverge.com/2026/3/10/article-slug",
  "title": "Article Title",
  "link": "https://www.theverge.com/2026/3/10/article-slug",
  "description": "A short summary of the article...",
  "published_at": "2026-03-10T14:30:00Z",
  "fetched_at": "2026-03-10T14:32:00Z"
}
```

---

## 6. Data Flow

```
┌────────────┐      JSON       ┌────────────┐      AMQP       ┌─────────────┐
│            │  ──────────────► │            │  ─────────────► │             │
│  /sources  │   parse YAML    │  Fetcher   │   publish msg   │  RabbitMQ   │
│  (*.yaml)  │                 │            │                 │             │
└────────────┘                 └────────────┘                 └──────┬──────┘
                                                                     │
                                                                     │ consume
                                                                     ▼
┌────────────┐     SELECT      ┌────────────┐     INSERT      ┌─────────────┐
│            │  ◄────────────  │            │  ◄────────────  │             │
│   Client   │                 │    API     │                 │  Processor  │
│            │  ──────────────►│  (read-only)│                │  (writer)   │
└────────────┘    HTTP/JSON    └─────┬──────┘                 └──────┬──────┘
                                     │                               │
                                     │          ┌──────────┐         │
                                     └─────────►│ Postgres │◄────────┘
                                      SELECT    └──────────┘  INSERT/UPDATE
```

---

## 7. Development & Testing Workflow

### Local Development

Provide a `docker-compose.yml` with RabbitMQ and Postgres:

```bash
docker compose up -d
```

### Testing Feeds

- Create a `/test` directory with sample `.xml` files
- Run a mock HTTP server serving those files
- Validate that the Fetcher correctly parses varied RSS/Atom formats

### Open-Source Contributor Flow

```
Fork  ──►  Add mysource.yaml  ──►  Submit PR  ──►  CI validates feed URL  ──►  Merge
```
