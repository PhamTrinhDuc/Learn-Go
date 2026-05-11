package redis

import (
	"iter"

	"github.com/redis/go-redis/v9"
	"google.golang.org/adk/session"
)

// redisEvents implements session.Events with live Redis reads.
// When filtered is true, the cached slice is the authoritative source (e.g.
// after Get applied NumRecentEvents / After filters) and loadFromRedis returns
// it directly without re-fetching.

type redisEvents struct {
	client   redis.Client
	key      string
	cached   []*session.Event
	isFilter bool
}

func newRedisEvents(events []*session.Event, key string, client *redis.Client) *redisEvents {
	if events == nil {
		events = make([]*session.Event, 0)
	}
	return &redisEvents{
		client: *client,
		key:    key,
		cached: events,
	}
}

func newRedisEventsWithFilter(events []*session.Event, key string, client *redis.Client) *redisEvents {
	if events == nil {
		events = make([]*session.Event, 0)
	}
	return &redisEvents{
		client:   *client,
		key:      key,
		cached:   events,
		isFilter: true,
	}
}

func (s *redisEvents) All() iter.Seq[*session.Event] {
	return nil
}

func (s *redisEvents) Len() int {
	return 0
}

func (s *redisEvents) At(i int) *session.Event {
	return nil
}

var _ session.Events = (*redisEvents)(nil)
