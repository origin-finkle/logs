package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/origin-finkle/logs/internal/logger"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
	"github.com/origin-finkle/logs/internal/wowhead"
)

var cifMutex sync.RWMutex

func GetCastInFight(ctx context.Context, id int64) (*models.CastInFight, error) {
	ctx = logger.ContextWithLogger(ctx, logger.FromContext(ctx).WithField("spell_id", id))
	cifMutex.RLock()
	if v, ok := data.CastInFight[id]; ok {
		cifMutex.RUnlock()
		return v, nil
	}
	cifMutex.RUnlock()
	spell, err := wowhead.GetSpell(ctx, id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Warnf("did not load item %d", id)
		return nil, ErrCastInFightNotFound
	}
	cif := &models.CastInFight{
		SpellID: spell.SpellID,
		Name:    spell.Name,
		Rank:    spell.Rank,
		Display: true,
	}
	if cif.Rank == 0 {
		cif.Rank = 1
	}
	switch true {
	// TODO: differenciate items used
	/*case cif.Rank > 0:
	cif.Type = "spell"*/
	default:
		cif.Type = "spell"
	}
	SetCastInFight(cif)
	postAddCIF(cif)
	return cif, nil
}

func markCIFAsInvalid(lowerRanked, higherRanked *models.CastInFight) {
	if lowerRanked.Invalid {
		already := data.CastInFight[lowerRanked.SuggestedSpellID]
		if already.Rank > higherRanked.Rank {
			// do not update as the suggested spell has a higher rank than the current one
			return
		}
	}
	lowerRanked.Invalid = true
	lowerRanked.InvalidReason = remark.Type_CastHigherRankAvailable
	lowerRanked.SuggestedSpellID = higherRanked.SpellID
}

func postAddCIF(cif *models.CastInFight) {
	cifMutex.RLock()
	defer cifMutex.RUnlock()
	// find similar cif, if rank < current one, then mark as invalid
	for _, castInFight := range data.CastInFight {
		if castInFight.Name == cif.Name {
			if castInFight.Rank < cif.Rank {
				markCIFAsInvalid(castInFight, cif)
			} else if castInFight.Rank > cif.Rank {
				markCIFAsInvalid(cif, castInFight)
			}
		}
	}
}

func SetCastInFight(v *models.CastInFight) {
	cifMutex.Lock()
	defer cifMutex.Unlock()
	v.TextRule = v.CommonConfig.String()
	data.CastInFight[v.SpellID] = v
}

var (
	ErrCastInFightNotFound = fmt.Errorf("cast in fight not found")
)
