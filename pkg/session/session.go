package session

import (
	"encoding/json"
	"time"

	"github.com/mirzakhany/pm/pkg/kv"
)

// Get get a data from session if its available
func Get(key string, data interface{}) error {
	val, err := kv.Get().GetString(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), data)
}

// Set set new key/value into session data
func Set(key string, data interface{}, duration time.Duration) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return kv.Get().Set(key, string(val), duration)
}
