package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUpRedisState() redisState {
	return redisState{
		service: redisSvc,
		client:  redisSvc.client,
		key:     redisSvc.appStateKey(appName),
		ttl:     redisSvc.ttl,
		appName: appName,
		userID:  userID,
		data: map[string]any{
			"userName": userName,
			"userId":   userID,
		},
	}
}

func TestMethodsState(t *testing.T) {
	ctx := context.Background()
	redisState := setUpRedisState()

	t.Run(getName("check data state"), func(t *testing.T) {
		assert.Equal(t, redisState.data, map[string]any{
			"userName": userName,
			"userId":   userID,
		})
	})

	t.Run(getName("check app state"), func(t *testing.T) {
		keyApp := redisState.service.appStateKey(appName)
		err := redisState.Set(
			keyApp,
			map[string]any{
				"app:theme": "dark",
				"app:lang":  "vi",
			},
		)
		assert.NoError(t, err)

		data, err := redisState.client.HGet(ctx, keyApp, "app:theme").Result()
		assert.NoError(t, err)

		assert.Equal(t, map[string]any{
			"app:theme": "dark",
			"app:lang":  "vi",
		}, data)

		redisState.service.FlushDB(ctx)
	})
}
