# 02-Requirements & Capacity Planning

## 1. Overview

This document defines the **functional requirements**, **non-functional requirements**, and **capacity estimations** for a scalable real-time auction / bidding platform.

The system must handle **massive last-second bidding spikes** while guaranteeing **fairness and correctness**.

---

# 2. Functional Requirements

## 2.1 User Management

Users must be able to:

- Register and login
- View auctions
- Place bids
- Receive notifications
- View bid history
- View auctions they created or participated in

---

## 2.2 Auction Lifecycle

Sellers must be able to:

| Action | Description |
|---|---|
| Create Auction | Create an auction with start & end time |
| Start Auction | Automatically start at scheduled time |
| End Auction | Automatically end at scheduled time |
| Cancel Auction | Allowed before first bid |
| View Auction | Publicly viewable auction page |
| Determine Winner | Highest valid bid at end time |

Auction states:
    SCHEDULED → ACTIVE → ENDED → SETTLED

---

## 2.3 Bidding

Users must be able to:

- Place a bid on an active auction
- See current highest bid
- See bid history

A bid must be rejected if:
- Auction is not active
- Bid amount ≤ current highest bid
- User is the auction owner

Rules:

- Bids are append-only (never update or delete)
- Highest bid wins
- If two bids have same value → earliest wins

---

## 2.4 Real-Time Updates

Users watching an auction must receive:

- New bid updates in real time
- "You were outbid" notification
- Auction ending reminder
- Auction winner announcement

Target latency:

- **99th percentile: < 200ms**

---

# 3. Non-Functional Requirements

## 3.1 Availability

- Target uptime: **99.9%**
- Auction closing must be highly reliable

---

## 3.2 Consistency

Bidding requires **strong consistency**:

We must guarantee:
- No lost bids
- No double winners
- Correct final price

Consistency model:

- **Strong consistency** for bidding operations
- **Eventual consistency** for reads & analytics

---

## 3.3 Durability

- All bids must be durably stored
- No data loss during failures
- Recovery within minutes

---

## 3.4 Scalability

System must scale to handle:

- **10,000+ concurrent auctions**
- **1,000+ bids per second**
- **100,000+ daily active users**
- **Spiky traffic patterns**

---

## 3.5 Latency

| Operation | Target Latency |
|---|---|
| Place bid | **< 200ms** |
| Get auction | **< 150ms** |
| Real-time update | **< 1 sec (99th percentile)** |

---

# 4. Capacity Estimation

## 4.1 Assumptions

| Metric | Value |
|---|---|
| Registered users | 50 million |
| Daily active users | 5 million |
| Concurrent users | 1 million |
| Active auctions | 10 million |
| Avg auction duration | 24 hours |
---

---

## 4.2 Auction Read Traffic

Assume each active user views 20 auctions/day.

Daily auction reads = 5M × 20 = 100M reads/day
Reads/sec ≈ 1,200 RPS
Peak ≈ 5,000 RPS

---

## 4.3 Bid Traffic (Critical Path)

Assume:
- 10% of users bid daily → 500k bidders/day
- Avg bids per bidder = 10
Daily bids = 5M bids/day
Average ≈ 58 bids/sec

Due to last-second spikes:
Peak bids/sec ≈ 50,000 bids/sec


This is the most important scaling number in the system.

---

## 4.4 Storage Estimation

### Auctions

Auction size ≈ 1 KB
10M active auctions → 10 GB
Historical (5 years) → ≈ 18 TB


---

### Bids

Bid size ≈ 200 bytes
5M bids/day × 365 × 5 years ≈ 9.1B bids
Storage ≈ 1.8 TB


Bids are append-only → ideal for event storage.

---

## 4.5 Real-Time Event Bandwidth

Assume:
- 1,000 watchers per hot auction
- 50k bids/sec peak
- 200 bytes per event
Outbound traffic ≈ 10 MB/sec (~80 Mbps)


This requires event streaming & pub/sub.

---

# 5. Key Technical Challenges

## Challenge 1 — Concurrent Bidding
Two users may bid simultaneously.

We must prevent:
- Lost updates
- Race conditions
- Double winners

---

## Challenge 2 — Auction Closing Spike
Traffic may spike 1000× near auction deadlines.

We must handle:
- Hot auctions
- Sudden traffic bursts
- Database contention

---

## Challenge 3 — Real-Time Fan-out
One bid must notify thousands of watchers.

This requires:
- Pub/Sub
- WebSockets
- Event streaming

---

# 6. API Preview

### Auction APIs

POST /auctions
GET /auctions/{id}
GET /auctions?status=active


### Bidding APIs


POST /auctions/{id}/bids
GET /auctions/{id}/bids


---

## Summary

The system must support:

- Millions of auctions
- Tens of thousands of bids/sec
- Real-time updates
- Strong consistency
- Horizontal scalability