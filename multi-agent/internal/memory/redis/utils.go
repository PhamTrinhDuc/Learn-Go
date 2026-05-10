package redis

import (
	"encoding/json"
	"strings"

	"google.golang.org/adk/session"
)

// extractStateDeltas separate data state of App, User, Session follow prefix of key
func extractStateDeltas(state map[string]any) (appDelta, userDelta, sessionDelta map[string]any) {

	if state == nil {
		return appDelta, userDelta, sessionDelta
	}

	appDelta = make(map[string]any)
	userDelta = make(map[string]any)
	sessionDelta = make(map[string]any)

	for key, value := range state {
		// {"app:theme": "dark"} => {"theme": "dark"}
		if cleanedKey, found := strings.CutPrefix(key, session.KeyPrefixApp); found {
			appDelta[cleanedKey] = value
			// {"user:name": "Jiyuu"} => {"name": "Jiyuu"}
		} else if cleanedKey, found := strings.CutPrefix(key, session.KeyPrefixUser); found {
			userDelta[cleanedKey] = value
		} else {
			// {"temp:token_count": 1234} => {"token_count": 1234}
			cleanedKey, found := strings.CutPrefix(key, session.KeyPrefixTemp)
			if found {
				sessionDelta[cleanedKey] = value
			}
		}
	}
	return appDelta, userDelta, sessionDelta
}

// unmarshalHashFields converts a Redis HASH result back to map[string]any by
// JSON-decoding each value.
func unmarshalHashFields(data map[string]string) map[string]any {
	result := make(map[string]any, len(data))
	for key, value := range data {
		var parsedValue any
		if err := json.Unmarshal([]byte(value), &parsedValue); err != nil {
			result[key] = value
		}
		result[key] = parsedValue
	}
	return result
}

// marshalHashFields converts a map[string]any to a map[string]string suitable
// for Redis HASH storage by JSON-encoding each value.
func marshalHashFields(data map[string]any) map[string]string {
	result := make(map[string]string, len(data))
	for key, value := range data {
		data, err := json.Marshal(value)
		if err != nil {
			continue
		}
		result[key] = string(data)
	}
	return result
}
