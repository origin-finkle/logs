package wowhead

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/sirupsen/logrus"
)

func getWowheadItem(ctx context.Context, itemID int64) (*WowheadItem, error) {
	retry := retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)
	var resp *http.Response
	if err := retry.Run(func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://fr.tbc.wowhead.com/item=%d&xml", itemID), nil)
		if err != nil {
			return err
		}
		resp, err = Client.Do(req)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("query to wowhead returned %d", resp.StatusCode)
	}
	v, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var item WowheadItem
	err = xml.Unmarshal(v, &item)
	if err != nil {
		return nil, err
	}
	if item.Item.Name == "" {
		var wErr Error
		err = xml.Unmarshal(v, &wErr)
		if err != nil {
			return nil, err
		}
		logrus.WithError(wErr).Warnf("item %d not found", itemID)
		return nil, wErr
	}
	return &item, nil
}

func GetItem(ctx context.Context, itemID int64) (*Item, error) {
	item, err := getWowheadItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	jsonEquip := "{" + item.Item.JsonEquip + "}"
	var iJSONEquip itemJSONEquip
	if err := json.Unmarshal([]byte(jsonEquip), &iJSONEquip); err != nil {
		return nil, err
	}
	return &Item{
		Name:    item.Item.Name,
		Slot:    item.Item.InventorySlot.Text,
		Sockets: iJSONEquip.Sockets,
	}, nil
}

type itemJSONEquip struct {
	Sockets int64 `json:"nsockets"`
}

type Item struct {
	Name    string `json:"name"`
	Slot    string `json:"slot"`
	Sockets int64  `json:"sockets"`
}
