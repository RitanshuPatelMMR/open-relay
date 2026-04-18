package worker

import (
        "context"
        "encoding/json"
        "log"
        "time"

        "github.com/redis/go-redis/v9"
        "github.com/ritanshupatel/openrelay/worker-service/deadletter"
        "github.com/ritanshupatel/openrelay/worker-service/store"
)

type Pool struct {
	redis       *redis.Client
	events      *store.EventStore
	deadletter  *deadletter.Handler
	workerCount int
	maxRetries  int
	timeoutSecs int
}

func NewPool(
	redisClient *redis.Client,
	events *store.EventStore,
	dl *deadletter.Handler,
	workerCount, maxRetries, timeoutSecs int,
) *Pool {
	return &Pool{
		redis:       redisClient,
		events:      events,
		deadletter:  dl,
		workerCount: workerCount,
		maxRetries:  maxRetries,
		timeoutSecs: timeoutSecs,
	}
}

func (p *Pool) Start(ctx context.Context) {
	p.redis.XGroupCreateMkStream(ctx, "events", "workers", "0")
	log.Printf("starting %d workers", p.workerCount)
	for i := 0; i < p.workerCount; i++ {
		workerID := i
		go p.runWorker(ctx, workerID)
	}
}

func (p *Pool) runWorker(ctx context.Context, id int) {
	log.Printf("worker %d started", id)
	consumerName := "worker-" + string(rune('0'+id))

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopped", id)
			return
		default:
		}

		streams, err := p.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "workers",
			Consumer: consumerName,
			Streams:  []string{"events", ">"},
			Count:    1,
			Block:    2 * time.Second,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("worker %d redis read error: %v", id, err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				eventID, ok := msg.Values["event_id"].(string)
				if !ok {
					p.redis.XAck(ctx, "events", "workers", msg.ID)
					continue
				}
				p.processEvent(ctx, eventID, msg.ID)
			}
		}
	}
}

func (p *Pool) processEvent(ctx context.Context, eventID, msgID string) {
	log.Printf("processing event: %s", eventID)

	event, endpoint, err := p.events.GetEventWithEndpoint(ctx, eventID)
	if err != nil {
		log.Printf("failed to fetch event %s: %v", eventID, err)
		p.redis.XAck(ctx, "events", "workers", msgID)
		return
	}

	if endpoint.DestinationURL == "" {
		log.Printf("event %s has no destination URL, skipping", eventID)
		p.redis.XAck(ctx, "events", "workers", msgID)
		return
	}

	for attempt := 1; attempt <= p.maxRetries; attempt++ {
		result := Deliver(ctx, endpoint.DestinationURL, event.Method, event.Headers, event.Payload, p.timeoutSecs)

		if result.Err == nil && result.StatusCode >= 200 && result.StatusCode < 300 {
			p.events.MarkDelivered(ctx, eventID)
                    p.events.LogAttempt(ctx, eventID, result.StatusCode, result.Body, result.DurationMs, "", nil)
                    p.redis.XAck(ctx, "events", "workers", msgID)
                    statusUpdate, _ := json.Marshal(map[string]string{
                            "type":      "delivery_update",
                            "event_id":  eventID,
                            "status":    "delivered",
                            "method":    event.Method,
                            "path":      "/",
                            "timestamp": time.Now().UTC().Format(time.RFC3339),
                    })
                    pubResult := p.redis.Publish(ctx, "delivery_updates", statusUpdate)
                    if pubResult.Err() != nil {
                            log.Printf("❌ redis publish failed: %v", pubResult.Err())
                    } else {
                            log.Printf("📡 published delivery_update for event %s (receivers: %d)", eventID, pubResult.Val())
                    }
                    log.Printf("✅ delivered event %s → %d in %dms", eventID, result.StatusCode, result.DurationMs)
                    return
		}

		errMsg := ""
		if result.Err != nil {
			errMsg = result.Err.Error()
		}

		if attempt < p.maxRetries {
			delay := RetryDelay(attempt)
			nextRetry := time.Now().Add(delay)
			p.events.LogAttempt(ctx, eventID, result.StatusCode, result.Body, result.DurationMs, errMsg, &nextRetry)
			p.events.IncrementAttempt(ctx, eventID)
			log.Printf("⚠️  event %s attempt %d failed (status %d), retrying in %s", eventID, attempt, result.StatusCode, delay)
			time.Sleep(delay)
		} else {
		p.events.MarkDelivered(ctx, eventID)
p.events.LogAttempt(ctx, eventID, result.StatusCode, result.Body, result.DurationMs, "", nil)
p.redis.XAck(ctx, "events", "workers", msgID)

statusUpdate, _ := json.Marshal(map[string]string{
        "type":      "delivery_update",
        "event_id":  eventID,
        "status":    "delivered",
        "method":    event.Method,
        "path":      "/",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
})
pubResult := p.redis.Publish(ctx, "delivery_updates", statusUpdate)
if pubResult.Err() != nil {
        log.Printf("❌ redis publish failed: %v", pubResult.Err())
} else {
        log.Printf("📡 published delivery_update for event %s (receivers: %d)", eventID, pubResult.Val())
}
log.Printf("✅ delivered event %s → %d in %dms", eventID, result.StatusCode, result.DurationMs)
return
		}
	}
}