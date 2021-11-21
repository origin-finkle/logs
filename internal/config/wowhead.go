package config

import (
	"context"
	"strconv"
	"sync"

	"github.com/origin-finkle/logs/internal/wowhead"
)

type Wowhead struct {
	mu    sync.RWMutex             `json:"-"`
	Items map[string]*wowhead.Item `json:"items"`
}

func GetWowheadItem(ctx context.Context, itemID int64) (*wowhead.Item, error) {
	data.Wowhead.mu.RLock()
	if item, ok := data.Wowhead.Items[strconv.FormatInt(itemID, 10)]; ok {
		data.Wowhead.mu.RUnlock()
		return item, nil
	}
	data.Wowhead.mu.RUnlock()

	item, err := wowhead.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	data.Wowhead.mu.Lock()
	data.Wowhead.Items[strconv.FormatInt(itemID, 10)] = item
	data.Wowhead.mu.Unlock()

	return item, nil
}
