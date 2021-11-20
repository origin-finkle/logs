package analyze

import (
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/remark"
)

func checkTalents(fa *models.FightAnalysis) {
	if fa.Talents == nil {
		return
	}
	used := fa.Talents.Points[0] + fa.Talents.Points[1] + fa.Talents.Points[2]
	if used != 61 && used > 0 {
		fa.AddRemark(remark.InvalidTalentPoints{
			ExpectedPoints: 61,
			PointsUsed:     used,
		})
	}
}
