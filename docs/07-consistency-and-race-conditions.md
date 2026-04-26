# 07-Consistency and Race Conditions
📄 docs/07-consistency-and-race-conditions.md
# Consistency, Race Conditions & Edge Cases

## 1. Why this chapter matters

Distributed systems fail in subtle ways.

A bidding system must **never produce an incorrect winner**.

This chapter explores the hardest production scenarios and how we handle them.

---

# 2. Consistency Model Recap

We use **different consistency levels** for different paths:

| Component | Consistency |
|---|---|
| Redis bidding path | Strong (atomic Lua) |
| PostgreSQL storage | Eventual |
| Real-time updates | Best effort |

This hybrid model gives both **speed and correctness**.

---

# 3. Double Bidding Race Condition

Scenario:
Two bids arrive at the same millisecond.

Bid A → $100  
Bid B → $120

Redis Lua script ensures atomic execution:


Check price → Validate → Update → Publish event


Only one script runs at a time per auction key.

Result:
No lost updates, deterministic ordering.

---

# 4. Duplicate Requests (Retries)

Clients retry requests when:
- Network timeout
- Mobile connection drop
- Server slow response

This can cause **duplicate bids**.

Solution: Idempotency keys.

Each bid request includes:


idempotency_key (UUID)


Stored with the bid event.

If the same key appears again:
- Return previous result
- Do NOT create new bid

---

# 5. Out-of-Order Events

Workers process queue messages asynchronously.

Events may arrive out of order.

Example:

Bid #102 stored before Bid #101 ❌

Solution:
Use event timestamps + ordering when querying.

Winner query:


ORDER BY amount DESC, timestamp ASC


The correct winner is always deterministic.

---

# 6. Redis Crash Scenario

If Redis crashes we lose hot auction state.

Recovery strategy:

1) Load active auctions from PostgreSQL
2) Rebuild Redis cache
3) Resume bidding

Because PostgreSQL is source of truth, no data is lost.

---

# 7. Queue Failure Scenario

If workers stop consuming:

- Queue stores events durably
- Processing resumes when workers restart

This guarantees **no bid loss**.

---

# 8. Database Downtime

If PostgreSQL becomes unavailable:

- Redis still accepts bids
- Queue stores events
- Workers retry until DB recovers

System continues operating in **degraded mode**.

---

# 9. Network Partition Scenario

Hard distributed systems problem.

Example:
Realtime service cannot receive Pub/Sub events.

Impact:
Users may miss live updates.

Mitigation:
Clients periodically refresh auction state via REST.

System eventually becomes consistent.

---

# 10. Auction Ending Edge Case

What if a bid arrives exactly at the deadline?

We define strict rule:

A bid is valid if:

bid_timestamp <= auction_end_time


Redis uses server time to enforce fairness.

---

# 11. Clock Synchronization

Multiple servers must share consistent time.

We rely on:
- NTP synchronization
- Server-side timestamps (never trust client time)

---

# 12. Exactly-Once vs At-Least-Once

Queue processing is **at-least-once**.

This means events may be processed multiple times.

We guarantee correctness using:

- Idempotency keys
- Unique DB constraints

System behaves like **exactly-once**.

---

# 13. Split Brain Auctions (Future)

In multi-region deployments:
Two regions may accept bids simultaneously.

Solution (future chapter):

- Leader election per auction
- Region ownership
- Kafka partitioning

We prepare architecture for this evolution.

---

# 14. Summary

We designed safeguards for:

- Race conditions
- Duplicate requests
- Out-of-order events
- Crashes & downtime
- Network failures

The system is now **resilient and production-ready**.

 