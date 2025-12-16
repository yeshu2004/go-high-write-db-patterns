# High-Write SQL Insert Strategies in Go

This repository explores **different SQL insert strategies in Go** and compares their **performance characteristics** under high-write workloads — a very common scenario in **web scrapers, BFS crawlers, event pipelines, and ingestion services**.

The goal is simple:
> **Given N links in memory, what is the fastest and safest way to persist them into a MySQL database using Go?**

---

## Use Case

- URLs are collected in memory
- After every **500–1000 URLs**, data must be flushed to the database
- Requirements:
  - High throughput
  - Minimal DB overhead
  - No data loss OR partial
  - Predictable performance

---

## Strategies Compared

The following approaches were implemented and tested:

| Strategy | Description |
|-------|------------|
| **Single Insert per Row** | One `INSERT` per row using `db.Exec` |
| **Prepared Statements** | Prepare once, execute many times |
| **Batch Insert** | Single `INSERT` with multiple `VALUES` |
| **Transaction Inserts** | Multiple inserts wrapped in a transaction |
| **Transaction + Batch Insert** | Batch insert wrapped in a transaction |

---

## Performance Results

> Results measured on local MySQL with default config

### N = 1,000 rows
| Strategy | Time |
|-------|------|
| Single Insert | Very slow |
| Prepared Statements | Slow |
| Batch Insert | ~0.78–0.85s |
| Transaction Inserts | ~0.89–0.93s |
| **Transaction + Batch** | **~0.76–0.84s** |

### N = 10,000 rows
| Strategy | Time |
|-------|------|
| Single Insert | ~18s |
| Prepared Statements | ~17–21s |
| Batch Insert | ~1.15s |
| Transaction Inserts | ~2.6–2.7s |
| **Transaction + Batch** | **~1.22–1.55s** |

---

## Key Takeaways

- Never insert rows one-by-one in high-write workloads
- Prepared statements alone are not enough
- Batch inserts drastically reduce DB overhead
- Transactions reduce commit cost
- Batch + Transaction is the optimal strategy

---

## 🧪 Testing & Benchmarks

### Run tests
```bash
go test -v
```

### Run benchmarks
```bash
go test -bench=. -run=^$ -benchmem
```

Benchmarks **rollback transactions** to avoid polluting the database.

---

