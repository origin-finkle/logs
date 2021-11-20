package events_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/models"
	"github.com/origin-finkle/logs/internal/models/events"
)

func TestProcessEvent_CombatantInfo(t *testing.T) {
	ev := json.RawMessage(`{
		"timestamp": 312518,
		"type": "combatantinfo",
		"sourceID": 9,
		"gear": [
			{
				"id": 29086,
				"quality": 4,
				"icon": "inv_helmet_81.jpg",
				"itemLevel": 120,
				"permanentEnchant": 3001,
				"gems": [
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					},
					{
						"id": 25897,
						"itemLevel": 70,
						"icon": "inv_misc_gem_diamond_06.jpg"
					}
				]
			},
			{
				"id": 28731,
				"quality": 4,
				"icon": "inv_jewelry_necklace_32.jpg",
				"itemLevel": 115
			},
			{
				"id": 29089,
				"quality": 4,
				"icon": "inv_shoulder_44.jpg",
				"itemLevel": 120,
				"permanentEnchant": 2980,
				"gems": [
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					},
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					}
				]
			},
			{
				"id": 0,
				"quality": 1,
				"icon": "inv_axe_02.jpg",
				"itemLevel": 0
			},
			{
				"id": 28600,
				"quality": 4,
				"icon": "inv_chest_leather_07.jpg",
				"itemLevel": 115,
				"permanentEnchant": 1144
			},
			{
				"id": 28655,
				"quality": 4,
				"icon": "inv_belt_22.jpg",
				"itemLevel": 115
			},
			{
				"id": 28591,
				"quality": 4,
				"icon": "inv_pants_mail_15.jpg",
				"itemLevel": 115,
				"permanentEnchant": 2746,
				"gems": [
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					},
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					},
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					}
				]
			},
			{
				"id": 28752,
				"quality": 4,
				"icon": "inv_boots_chain_05.jpg",
				"itemLevel": 115,
				"permanentEnchant": 464,
				"gems": [
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					},
					{
						"id": 24029,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_livingruby_03.jpg"
					}
				]
			},
			{
				"id": 29249,
				"quality": 4,
				"icon": "inv_bracer_13.jpg",
				"itemLevel": 110,
				"permanentEnchant": 2617
			},
			{
				"id": 28521,
				"quality": 4,
				"icon": "inv_gauntlets_25.jpg",
				"itemLevel": 115,
				"permanentEnchant": 2322,
				"gems": [
					{
						"id": 30603,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_nightseye_03.jpg"
					},
					{
						"id": 30547,
						"itemLevel": 70,
						"icon": "inv_jewelcrafting_nobletopaz_03.jpg"
					}
				]
			},
			{
				"id": 29290,
				"quality": 4,
				"icon": "inv_jewelry_ring_62.jpg",
				"itemLevel": 130
			},
			{
				"id": 22939,
				"quality": 4,
				"icon": "inv_jewelry_ring_50naxxramas.jpg",
				"itemLevel": 83
			},
			{
				"id": 29376,
				"quality": 4,
				"icon": "inv_valentineperfumebottle.jpg",
				"itemLevel": 110
			},
			{
				"id": 30841,
				"quality": 3,
				"icon": "inv_misc_book_11.jpg",
				"itemLevel": 115
			},
			{
				"id": 28765,
				"quality": 4,
				"icon": "inv_misc_cape_06.jpg",
				"itemLevel": 125
			},
			{
				"id": 32451,
				"quality": 4,
				"icon": "inv_mace_47.jpg",
				"itemLevel": 123,
				"permanentEnchant": 2343
			},
			{
				"id": 29170,
				"quality": 4,
				"icon": "inv_misc_orb_01.jpg",
				"itemLevel": 105
			},
			{
				"id": 27886,
				"quality": 1,
				"icon": "spell_nature_natureresistancetotem.jpg",
				"itemLevel": 112
			},
			{
				"id": 0,
				"quality": 1,
				"icon": "inv_axe_02.jpg",
				"itemLevel": 0
			}
		],
		"auras": [
			{
				"source": 9,
				"ability": 33268,
				"stacks": 1,
				"icon": "spell_misc_food.jpg",
				"name": "Well Fed"
			},
			{
				"source": 10,
				"ability": 469,
				"stacks": 1,
				"icon": "ability_warrior_rallyingcry.jpg",
				"name": "Commanding Shout"
			},
			{
				"source": 9,
				"ability": 28491,
				"stacks": 1,
				"icon": "inv_potion_142.jpg",
				"name": "Healing Power"
			},
			{
				"source": 9,
				"ability": 39627,
				"stacks": 1,
				"icon": "inv_potion_155.jpg",
				"name": "Elixir of Draenic Wisdom"
			},
			{
				"source": 11,
				"ability": 27143,
				"stacks": 1,
				"icon": "spell_holy_greaterblessingofwisdom.jpg",
				"name": "Greater Blessing of Wisdom"
			},
			{
				"source": 4,
				"ability": 27127,
				"stacks": 1,
				"icon": "spell_holy_arcaneintellect.jpg",
				"name": "Arcane Brilliance"
			},
			{
				"source": 5,
				"ability": 25898,
				"stacks": 1,
				"icon": "spell_magic_greaterblessingofkings.jpg",
				"name": "Greater Blessing of Kings"
			},
			{
				"source": 5,
				"ability": 27149,
				"stacks": 1,
				"icon": "spell_holy_devotionaura.jpg",
				"name": "Devotion Aura"
			}
		],
		"expansion": "tbc",
		"faction": 1,
		"specID": 0,
		"strength": 96,
		"agility": 99,
		"stamina": 570,
		"intellect": 592,
		"spirit": 531,
		"dodge": 0,
		"parry": 0,
		"block": 0,
		"armor": 3804,
		"critMelee": 0,
		"critRanged": 0,
		"critSpell": 0,
		"hasteMelee": 0,
		"hasteRanged": 0,
		"hasteSpell": 0,
		"hitMelee": 0,
		"hitRanged": 0,
		"hitSpell": 0,
		"expertise": 0,
		"resilienceCritTaken": 18,
		"resilienceDamageTaken": 18,
		"talents": [
			{
				"id": 8,
				"icon": "inv_axe_02.jpg"
			},
			{
				"id": 11,
				"icon": "spell_frost_frostbolt02.jpg"
			},
			{
				"id": 42,
				"icon": "inv_axe_02.jpg"
			}
		],
		"pvpTalents": [],
		"customPowerSet": [],
		"secondaryCustomPowerSet": [],
		"tertiaryCustomPowerSet": []
	}`)
	analysis := &models.Analysis{
		Data: make(map[int64]*models.PlayerAnalysis),
	}
	analysis.SetPlayerAnalysis(9, &models.PlayerAnalysis{
		Fights: make(map[string]*models.FightAnalysis),
	})
	pa := analysis.GetPlayerAnalysis(9)
	pa.SetFight(&models.FightAnalysis{Name: "test"})
	fa := pa.GetFight("test")
	err := events.Process(context.TODO(), ev, analysis, "test")
	td.CmpNoError(t, err)
	td.Cmp(t, fa.Talents.Points, [3]int64{8, 11, 42})
}
