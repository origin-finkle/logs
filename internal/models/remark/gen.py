from os import system

ItemWowheadAttr = {
    "name": "ItemWowheadAttr",
    "json": "item_wowhead_attr",
    "type": "string",
}
WowheadAttr = {"name": "WowheadAttr", "json": "wowhead_attr", "type": "string"}
Slot = {"name": "Slot", "json": "slot", "type": "string"}
EnchantID = {"name": "EnchantID", "json": "enchant_id", "type": "int64"}
Count = {"name": "Count", "type": "int", "json": "count"}
ExpectedPoints = {"name": "ExpectedPoints", "type": "int64", "json": "expected_points"}
PointsUsed = {"name": "PointsUsed", "type": "int64", "json": "points_used"}
SpellID = {"name": "SpellID", "type": "int64", "json": "spell_id"}
SuggestedSpellID = {
    "name": "SuggestedSpellID",
    "type": "int64",
    "json": "suggested_spell_id",
}
SpellWowheadAttr = {
    "name": "SpellWowheadAttr",
    "json": "spell_wowhead_attr",
    "type": "string",
}
HigherRankedSpellWowheadAttr = {
    "name": "HigherRankedSpellWowheadAttr",
    "json": "higher_ranked_spell_wowhead_attr",
    "type": "string",
}
PossibleCasts = {"name": "PossibleCasts", "json": "possible_casts", "type": "int64"}
ActualCasts = {"name": "ActualCasts", "json": "actual_casts", "type": "int64"}
Threshold = {"name": "Threshold", "json": "threshold", "type": "int64"}
ActualPercentageOfUse = {
    "name": "ActualPercentageOfUse",
    "json": "actual_percentage_of_use",
    "type": "int64",
}

db = {
    "InvalidGem": {
        "str": "invalid_gem",
        "fields": [
            ItemWowheadAttr,
            WowheadAttr,
        ],
    },
    "MissingGems": {
        "str": "missing_gems",
        "fields": [
            ItemWowheadAttr,
            Count,
        ],
    },
    "NoEnchant": {"str": "no_enchant", "fields": [ItemWowheadAttr, Slot]},
    "InvalidEnchant": {
        "str": "invalid_enchant",
        "fields": [
            ItemWowheadAttr,
            WowheadAttr,
            Slot,
            EnchantID,
        ],
    },
    "NoTemporaryEnchant": {"str": "no_temporary_enchant", "fields": [ItemWowheadAttr]},
    "InvalidTemporaryEnchant": {
        "str": "invalid_temporary_enchant",
        "fields": [ItemWowheadAttr],
    },
    "NoTemporaryEnchantButWindfury": {
        "str": "no_temporary_enchant_but_windfury",
        "fields": [ItemWowheadAttr],
    },
    "MissingBattleElixir": {
        "str": "missing_battle_elixir",
        "fields": [],
    },
    "MissingGuardianElixir": {
        "str": "missing_guardian_elixir",
        "fields": [],
    },
    "MissingFood": {
        "str": "missing_food",
        "fields": [],
    },
    "InvalidBattleElixir": {
        "str": "invalid_battle_elixir",
        "fields": [WowheadAttr],
    },
    "InvalidGuardianElixir": {
        "str": "invalid_guardian_elixir",
        "fields": [WowheadAttr],
    },
    "InvalidFood": {
        "str": "invalid_food",
        "fields": [WowheadAttr],
    },
    "InvalidTalentPoints": {
        "str": "invalid_talent_points",
        "fields": [ExpectedPoints, PointsUsed],
    },
    "CastHigherRankAvailable": {
        "str": "cast_higher_rank_available",
        "fields": [
            SpellID,
            SuggestedSpellID,
            SpellWowheadAttr,
            HigherRankedSpellWowheadAttr,
            Count,
        ],
    },
    "MissingItemInSlot": {"str": "missing_item_in_slot", "fields": [Slot]},
    "MetaNotActivated": {
        "str": "meta_not_activated",
        "fields": [ItemWowheadAttr, WowheadAttr],
    },
    "CouldMaximizeCasts": {
        "str": "could_maximize_casts",
        "fields": [
            SpellID,
            PossibleCasts,
            ActualCasts,
            Threshold,
            ActualPercentageOfUse,
            SpellWowheadAttr,
        ],
    },
}

all_fields = {}

types_file = "./internal/models/remark/gen_types.go"
with open(types_file, "w+") as ft:
    ft.write(
        """
// Code generated by gen.py; DO NOT EDIT
package remark

type Type string

const (
    """
    )
    for object_name, data in db.items():
        ft.write(f"""\tType_{object_name} = "{data['str']}"\n""")
        test_file = f"./internal/models/remark/gen_{data['str']}_test.go"
        file = f"./internal/models/remark/gen_{data['str']}.go"
        with open(file, "w+") as f:
            f.write(
                f"""// Code generated by gen.py; DO NOT EDIT
package remark

type {object_name} struct {{
"""
            )
            for field in data["fields"]:
                all_fields[field["name"]] = field
                f.write(
                    f"""\t{field['name']} {field['type']} `json:"{field['json']}"`\n"""
                )
            f.write(
                f"""}}
func ({object_name}) getType() Type {{ return Type_{object_name} }}
func ({object_name}) is(rt Type) bool {{ return rt == Type_{object_name} }}
func (md {object_name}) apply(r *Remark) {{
"""
            )
            for field in data["fields"]:
                f.write(f"\tr.{field['name']} = md.{field['name']}\n")
            f.write("}\n")
        with open(test_file, "w+") as f:
            f.write(
                f"""// Code generated by gen.py; DO NOT EDIT
package remark_test

import (
    "testing"

    "github.com/origin-finkle/logs/internal/models/remark"
)

func Test{object_name}_applyToRemark(t *testing.T) {{
    md := remark.{object_name}{{}}
    remark.New(md, "test")
}}
            """
            )
        system(f"gofmt -s -w {test_file}")
        system(f"gofmt -s -w {file}")
    ft.write(")\n")
system(f"gofmt -s -w {types_file}")

remark_file = "./internal/models/remark/gen_remark.go"
with open(remark_file, "w+") as f:
    f.write(
        """// Code generated by gen.py; DO NOT EDIT
package remark

type Remark struct {
	Type            Type   `json:"type"`
	Fight           string `json:"fight,omitempty"`
"""
    )
    for field in all_fields.values():
        f.write(
            f"""\t{field['name']} {field['type']} `json:"{field['json']},omitempty"`\n"""
        )
    f.write(
        """
	UUID            string `json:"uuid"`
} 
"""
    )
system(f"gofmt -s -w {remark_file}")
