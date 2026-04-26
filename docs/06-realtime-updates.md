# 06-Realtime Updates — WebSockets & Pub/Sub

## 1. Why real-time updates matter

Auction platforms are **live experiences**.

Users expect to instantly see:

- New bids appearing
- Price changes
- "You were outbid" alerts
- Auction ending countdown
- Winner announcement

Polling every few seconds is not acceptable.

Target latency:

< 1 second


This requires **push-based communication**.

---

# 2. Why HTTP Polling Is Not Enough

Naive approach:

Client requests auction state every 2 seconds.

Problems:

- Massive unnecessary traffic
- Slow updates
- Poor user experience
- Expensive at scale

Example:

1M users polling every 2 seconds →  
500k requests/sec 😱

We need persistent connections.

---

# 3. Solution — WebSockets

WebSockets allow:

- Persistent bidirectional connection
- Server can push updates instantly
- Low overhead after connection established

We introduce a new component:

Realtime Service (WebSocket Gateway)


Architecture update:


Clients ↔ Realtime Service (WebSockets)
↑
Pub/Sub
↑
Bidding Service


---

# 4. Realtime Service Responsibilities

The realtime service will:

- Manage WebSocket connections
- Subscribe users to auctions
- Broadcast bid updates
- Handle presence & disconnects

This service must handle **millions of connections**.

---

# 5. Auction Subscription Model

Users subscribe when they open an auction page.

Client sends:

SUBSCRIBE auction_id


Server keeps mapping:


auction_id → list of connected clients


This allows targeted broadcasting.

---

# 6. Broadcasting New Bids

When a bid is accepted:

1) Bidding service publishes event
2) Pub/Sub distributes event
3) Realtime service broadcasts to subscribers

Flow:


Bid Accepted → Publish Event → Broadcast → Clients update UI


---

# 7. Why We Need Pub/Sub

Realtime service may run on **many instances**.

Without Pub/Sub:
- Only users connected to same server get updates ❌

With Pub/Sub:
- All realtime servers receive events
- All connected users receive updates

We will use Redis Pub/Sub initially.

Later this can evolve to Kafka.

---

# 8. Event Types

Realtime service handles multiple event types.

### BidPlaced

{
type: "BID_PLACED",
auction_id: "123",
amount: 120
}


### Outbid Notification

{
type: "OUTBID",
auction_id: "123"
}


### AuctionEnded

{
type: "AUCTION_ENDED",
winner_id: "user_456"
}


---

# 9. Scaling WebSocket Connections

One server cannot hold millions of connections.

We scale horizontally:


Load Balancer
↓
Realtime Servers (many instances)


Key techniques:

- Stateless servers
- Sticky sessions (optional)
- Pub/Sub for cross-instance communication

---

# 10. Handling Disconnects & Reconnects

Connections will drop frequently.

Client must:
- Reconnect automatically
- Re-subscribe to auctions
- Fetch latest state via REST if needed

This ensures reliability.

---

# 11. Fan-out Problem

Popular auctions may have:


10k+ watchers per auction


One bid → thousands of messages.

Pub/Sub + WebSockets make this scalable.

---

# 12. Reliability Considerations

Realtime updates are **best effort**.

Source of truth remains:
- Redis (hot state)
- PostgreSQL (persistent state)

If a client misses events:
- Page refresh fetches latest data.

---

# 13. Summary

We added real-time capabilities using:

- WebSockets for persistent connections
- Pub/Sub for cross-instance broadcasting
- Subscription model for targeted updates

The system now supports **live auctions at scale**.

---
 