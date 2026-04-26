# 03-High-Level Architecture
# High-Level Architecture

## 1. Introduction

In this chapter we design the first **end-to-end architecture** of the bidding platform.

We will evolve the system gradually:

V1 → Monolith  
V2 → Scaled Monolith + Cache + Queue  
V3 → Microservices + Streaming  

This document focuses on **Version 1 (MVP architecture)** while preparing for future scaling.

---

# 2. Big Picture

At a high level the system has 4 main responsibilities:

1. API layer (HTTP)
2. Auction management
3. Bidding engine (critical path)
4. Real-time notifications

Initial architecture:
Clients → API → Application → Database


This is intentionally simple and production-realistic for early stages.

---

# 3. System Context Diagram
        ┌──────────────┐
        │   Web / App  │
        └──────┬───────┘
               │ HTTPS
               ▼
        ┌──────────────┐
        │   API Layer  │
        └──────┬───────┘
               │
               ▼
    ┌─────────────────────┐
    │   Auction Monolith   │
    │                     │
    │  - Auction Module   │
    │  - Bidding Module   │
    │  - Notification     │
    └──────┬──────────────┘
           │
           ▼
     ┌────────────┐
     │ PostgreSQL │
     └────────────┘


This architecture can realistically support **hundreds of bids/sec**.

And most importantly — it gives us a **simple base to evolve**.

---

# 4. Why Start With a Monolith?

Many beginners jump directly to microservices.  
Real companies do not.

Starting with a monolith gives:

### Benefits

- Faster development
- Simpler deployment
- Strong ACID transactions
- Easier debugging
- Lower infrastructure cost

A well-designed monolith can scale surprisingly far.

We will extract services **only when bottlenecks appear**.

---

# 5. Internal Architecture (Inside the Monolith)

We will use **Modular Monolith Architecture**.
auction-system/
├── cmd/api
├── internal/
│ ├── auction/
│ ├── bidding/
│ ├── user/
│ ├── notification/
│ ├── realtime/
│ └── platform/



Each module has its own:
- Handlers
- Services
- Repository
- Domain models

This keeps the codebase clean and ready for microservice extraction later.

---

# 6. Core Components

## 6.1 API Layer

Responsibilities:

- Authentication (later)
- Request validation
- Routing
- Rate limiting (later)

Initial APIs:
POST /auctions
GET /auctions/{id}
POST /auctions/{id}/bids
GET /auctions/{id}/bids



---

## 6.2 Auction Module

Responsible for:

- Creating auctions
- Updating auction status
- Ending auctions
- Fetching auction data

Auction state machine:
SCHEDULED → ACTIVE → ENDED → SETTLED



Auction status updates will later move to background workers.

---

## 6.3 Bidding Module (Critical Path 🔥)

This is the **heart of the system**.

Responsibilities:

- Validate bids
- Prevent race conditions
- Update highest bid atomically
- Store bid history

This path must be:
- Fast
- Correct
- Strongly consistent

In V1 we rely on **PostgreSQL transactions**.

---

## 6.4 Notification Module

Initial version uses:

- Email / Push simulation
- Later replaced by event-driven architecture

Triggers:
- New bid placed
- User outbid
- Auction ended

---

## 6.5 Realtime Module

Version 1:
- Simple WebSocket server

Used for:
- Live auction updates
- Broadcasting new bids

Later this will be replaced by Pub/Sub.

---

# 7. Database Design (Initial)

We use **PostgreSQL** because we need:

- ACID transactions
- Row-level locking
- Strong consistency

Initial tables:

### auctions
id (PK)
title
seller_id
starting_price
current_price
status
start_time
end_time
created_at

### bids

id (PK)
auction_id (FK)
user_id
amount
created_at


### users

id (PK)
name
email
created_at

Indexes will be added later.

---

# 8. Bidding Request Flow (End-to-End)

This is the most important flow.

User clicks "Place Bid"
↓
API receives request
↓
Validate auction is ACTIVE
↓
Start DB transaction
↓
Lock auction row (SELECT FOR UPDATE)
↓
Validate bid > current price
↓
Insert bid record
↓
Update auction current_price
↓
Commit transaction
↓
Broadcast real-time update


This guarantees **no race conditions** in V1.

---

# 9. Known Limitations of V1

This architecture will hit limits at:

- ~500–1000 bids/sec
- Real-time fan-out becomes expensive
- Hot auctions overload DB
- No background processing

These limitations are intentional.

They give us a reason to evolve the architecture in next chapters.

---

# 10. Evolution Plan

Next chapters will progressively introduce:

1. Redis (hot auctions cache)
2. Message queue (async processing)
3. Event streaming (Kafka)
4. Service decomposition
5. Horizontal scaling

We will transform the monolith into a **real production architecture** step-by-step.

