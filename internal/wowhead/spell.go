package wowhead

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/origin-finkle/logs/internal/logger"
)

type WowheadSpell struct {
	SpellID int64
	Name    string
	Rank    int64
}

func GetSpell(ctx context.Context, spellID int64) (*WowheadSpell, error) {
	retry := retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)
	var resp *http.Response
	if err := retry.Run(func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://fr.tbc.wowhead.com/spell=%d/", spellID), nil)
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
	spell := &WowheadSpell{SpellID: spellID}
	strStart := fmt.Sprintf(`g_spells[%d].tooltip_frfr = "`, spellID)
	idx := strings.Index(string(v), strStart)
	idx += len(strStart)
	nextCarriage := strings.Index(string(v[idx:]), "\";\n")
	str := string(v[idx : idx+nextCarriage])
	str = strings.Replace(str, `\/`, `/`, -1)
	str = strings.Replace(str, `\"`, `"`, -1)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		logger.FromContext(ctx).WithError(err).Warn("did not load document")
		return nil, err
	}
	name := doc.Find("b.whtt-name").Text()
	name, err = strconv.Unquote(`"` + name + `"`)
	if err != nil {
		return nil, err
	}
	spell.Name = name
	doc.Find(".q0").Each(func(idx int, s *goquery.Selection) {
		if idx > 0 {
			return // only the first match matters
		}
		html, err := s.Html()
		if err != nil {
			logger.FromContext(ctx).WithError(err).Warn("could not convert to html")
			return
		}
		html = strings.Replace(html, "<br>", "\n", -1)
		html = strings.Replace(html, "<br/>", "\n", -1)
		str := strings.Split(html, "\n")[0]
		if str == "Talent" || str == "Raciale" || str == "Changeforme" || str == "gnome Racial" || str == "Racial" || str == "Invocation" || str == "humain Racial" || str == "Ma\u00eetre" {
			return
		}
		for _, prefix := range []string{"Rang ", "Rg "} {
			rank, err := strconv.ParseInt(strings.TrimPrefix(str, prefix), 10, 64)
			if err != nil {
				logger.FromContext(ctx).WithError(err).Warnf("could not parse rank string '%s'", s.Text())
				continue
			}
			spell.Rank = rank
			break
		}

	})
	return spell, nil
}
