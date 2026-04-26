# 05-Scaling
# Scaling the Hot Path — Redis & Queues

## 1. Why we need to scale

Our V1 architecture uses PostgreSQL transactions and row locking.

This guarantees correctness, but eventually it hits limits.

The main bottleneck appears when **many users bid on the same popular auction**.

Symptoms we will observe in production:

- DB CPU spikes
- Lock wait times increase
- API latency increases
- Connection pool saturation
- Bid placement slows down

We must optimize the **hot path**:

Place Bid → Update Current Price


---

# 2. Identify the real bottleneck

Not all auctions are equal.

Typical production pattern:

| Auction Type | Traffic |
|---|---|
| 95% of auctions | Very low traffic |
| 4% of auctions | Moderate traffic |
| 1% of auctions | Extremely hot 🔥 |

That 1% generates most traffic.

These are called:

Hot Auctions


We must optimize specifically for them.

---

# 3. Scaling Strategy Overview

We introduce two new components:


Redis → Hot auction cache (fast reads/writes)
Queue → Smooth traffic spikes


New architecture:


Clients → API → Bidding Service
↓
Redis (hot path)
↓
Queue (async)
↓
PostgreSQL (source of truth)


PostgreSQL remains the **source of truth**, but is no longer in the immediate critical path.

---

# 4. Role of Redis in the Bidding Path

Redis will store **active auctions in memory**.

Why Redis?

- In-memory → extremely fast
- Atomic operations
- High throughput (100k+ ops/sec)
- Perfect for short-lived hot data

We store for active auctions:


auction:{id}

current_price
end_time
version

Redis becomes the **real-time bidding engine**.

---

# 5. New Bid Flow (Scaled Version)

### Step-by-step flow

User places bid
API validates request
Bid sent to Redis first
Redis atomically validates & updates price
Event pushed to queue
Worker persists bid to PostgreSQL

This turns the system into **write-behind architecture**.

Users get fast responses while DB is updated asynchronously.

---

# 6. Atomic Bidding in Redis

We must still prevent race conditions.

We use **Lua scripts** for atomic updates.

Why Lua?

Redis executes Lua scripts atomically.

This gives us:
- Validation
- Comparison
- Update

All in one atomic operation.

Pseudo logic:


IF auction not active → reject
IF bid <= current_price → reject
UPDATE current_price
RETURN success


This replaces DB row locking.

---

# 7. Why introduce a Queue?

Even Redis cannot absorb sudden spikes forever.

During auction ending:
- 50k bids/sec may arrive
- Database cannot write this fast

Queue acts as a **shock absorber**.

Benefits:

| Problem | Queue Solution |
|---|---|
| Traffic spikes | Smooth burst traffic |
| DB overload | Async writes |
| Retries | Safe reprocessing |
| Failures | Durable storage |

This is called:

Event-driven architecture


---

# 8. Write-Behind Pattern

We now separate responsibilities:

| Component | Responsibility |
|---|---|
| Redis | Real-time bidding |
| Queue | Buffering events |
| Workers | Persist to DB |
| PostgreSQL | Source of truth |

User gets response immediately after Redis success.

DB eventually becomes consistent.

This is **eventual consistency** outside the critical path.

---

# 9. Bid Event Structure

Each accepted bid produces an event:


BidPlacedEvent

auction_id
user_id
amount
timestamp
idempotency_key

Workers consume events and store them in PostgreSQL.

---

# 10. Handling Failures

### Redis crash
We rebuild state from PostgreSQL + event log.

### Worker crash
Queue retains events → retry later.

### DB downtime
Queue buffers writes until DB recovers.

System becomes **resilient and durable**.

---

# 11. Benefits of the New Architecture

| Metric | Before | After |
|---|---|---|
| Bid latency | ~150–300 ms | ~5–20 ms |
| Max throughput | ~1k bids/sec | 50k+ bids/sec |
| DB load | High | Reduced |
| Spike handling | Poor | Excellent |

This is a major scalability milestone.

---

# 12. New Tradeoffs Introduced

Scaling always introduces tradeoffs.

We now have:

- Eventual consistency in DB
- More infrastructure complexity
- Need for background workers
- Need for monitoring

These tradeoffs are acceptable for scale.

---

# 13. Summary

We transformed the system from:


DB-centric → Cache + Queue architecture


We introduced:

- Redis for hot auctions
- Queue for buffering
- Workers for persistence
- Event-driven design
 