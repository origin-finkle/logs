package wowhead

import (
	"encoding/xml"
	"net/http"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var (
	Client = &http.Client{
		Transport: transport,
	}
	transport = &limitedTransport{
		sem: semaphore.NewWeighted(8),
	}
)

type WowheadItem struct {
	XMLName xml.Name `xml:"wowhead"`
	Text    string   `xml:",chardata"`
	Item    struct {
		Text    string `xml:",chardata"`
		ID      string `xml:"id,attr"`
		Name    string `xml:"name"`
		Level   string `xml:"level"`
		Quality struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"quality"`
		InventorySlot struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"inventorySlot"`
		JSON      string `xml:"json"`
		JsonEquip string `xml:"jsonEquip"`
		Link      string `xml:"link"`
	} `xml:"item"`
}

type Error struct {
	Err string `xml:"error"`
}

func (e Error) Error() string {
	return e.Err
}

type limitedTransport struct {
	sem *semaphore.Weighted
}

func (t *limitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.sem.Acquire(req.Context(), 1); err != nil {
		return nil, err
	}
	defer t.sem.Release(1)
	logger.FromContext(req.Context()).WithFields(logrus.Fields{
		"wowhead_url": req.URL.String(),
	}).Debug("requesting wowhead")
	return http.DefaultTransport.RoundTrip(req)
}
