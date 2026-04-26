# 04-Deep Dive Bidding Engine

## 1. Why the bidding engine is the hardest part

The bidding engine is the **critical path** of the entire system.

If this part fails, the platform loses trust immediately.

The engine must guarantee:

- No lost bids
- No race conditions
- Deterministic winner selection
- Strong consistency
- Low latency under heavy concurrency

This chapter explains how we achieve this step-by-step.

---

# 2. The Core Problem — Concurrent Bids

Imagine an auction ending soon.

Two users submit bids at the exact same millisecond:

User A → $100  
User B → $110

If not handled correctly, the system could:

- Accept both bids incorrectly
- Lose one bid
- Select wrong winner
- Corrupt auction price

This is the classic **race condition problem**.

We must design the system so that **only one bid can update the auction at a time**.

---

# 3. The Naive Approach (Incorrect)

A beginner implementation might do:

1) Read current price  
2) Validate bid > current price  
3) Insert bid  
4) Update auction price  

This fails because two requests may read the same price simultaneously.

This creates **lost update problem**.

We need **atomic bidding**.

---

# 4. Correct Approach — Transaction + Row Locking

In Version 1 we rely on PostgreSQL transactions.

Key idea:
We lock the auction row while processing a bid.

PostgreSQL provides:
SELECT ... FOR UPDATE


This creates a **row-level lock**.

Only one transaction can hold this lock at a time.

All other bids must wait in queue.

This guarantees **sequential bid processing per auction**.

---

# 5. Bid Placement Algorithm

### Step-by-step flow
BEGIN TRANSACTION

SELECT auction WHERE id = ? FOR UPDATE
Validate auction status == ACTIVE
Validate current_time < end_time
Validate bid_amount > current_price
INSERT INTO bids
UPDATE auctions SET current_price = bid_amount

COMMIT


This ensures:
- No concurrent updates
- No inconsistent state
- Guaranteed ordering

---

# 6. Why This Works

Row locking gives us:

| Property | Result |
|---|---|
| Isolation | Only one bidder updates auction at a time |
| Atomicity | Bid + price update happen together |
| Durability | Stored safely in DB |
| Consistency | Auction always valid |

This is why we start with **PostgreSQL-first architecture**.

---

# 7. Handling High Contention Auctions

Popular auctions become **hot rows**.

Thousands of bids may queue on the same row lock.

This is actually OK for V1.

Why?

Because it guarantees fairness:
Bids are processed in the exact order they arrive.

Later we will optimize this using Redis and queues.

---

# 8. Preventing Last-Second Sniping

Many auction systems extend auction time if bids arrive near the end.

Rule:
If a bid arrives in the last 10 seconds → extend auction by 10 seconds.

Implementation inside transaction:
IF end_time - now < 10 seconds:
end_time += 10 seconds


This prevents unfair last-millisecond wins.

---

# 9. Idempotency (Very Important)

Network retries can send the same bid multiple times.

We must prevent duplicate bids.

Solution: **Idempotency Key**

Each bid request contains:
idempotency_key (UUID)


Database constraint:

UNIQUE(idempotency_key)


If duplicate key → reject bid.


If request retries:
- First succeeds
- Retries safely ignored

---

# 10. Bid Ordering Guarantee

If two bids have same amount:

Earliest timestamp wins.

We guarantee deterministic order using:
ORDER BY amount DESC, created_at ASC


Winner query:


SELECT * FROM bids
WHERE auction_id = ?
ORDER BY amount DESC, created_at ASC
LIMIT 1


This ensures:
- Highest bid wins
- Earliest bid wins ties
- No ambiguity

---

# 11. Failure Scenarios

### DB crash during transaction
Transaction rolls back automatically → no partial updates.

### App crash after commit
Bid already stored → safe.

### Network timeout after commit
Client retries → idempotency prevents duplicates.

System remains correct.

---

# 12. Current Limitations

Row locking works well until:

- Thousands of bids/sec on same auction
- DB becomes bottleneck
- Connection pool saturation
- Increased latency

This is expected.

This is where we evolve the architecture in the next chapter.

---

# 13. Summary

Version 1 bidding engine uses:

- PostgreSQL transactions
- Row-level locking
- Idempotency keys
- Deterministic ordering

This gives us a **correct and reliable foundation**.

Next we will learn how to scale this using **Redis + queues**.