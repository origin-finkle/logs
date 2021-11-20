package wowhead_test

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/wowhead"
)

func TestGetItem(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(wowhead.Client)
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("GET", "https://fr.tbc.wowhead.com/item=1&xml", httpmock.NewStringResponder(200, `<?xml version="1.0" encoding="UTF-8"?><wowhead><error>Item not found!</error></wowhead>`))
	_, err := wowhead.GetItem(context.TODO(), 1)
	td.CmpError(t, err)

	httpmock.RegisterResponder("GET", "https://fr.tbc.wowhead.com/item=28514&xml", httpmock.NewStringResponder(200, `<?xml version="1.0" encoding="UTF-8"?><wowhead><item id="28514"><name><![CDATA[Brassards de malignité]]></name><level>115</level><quality id="4">Épique</quality><class id="4"><![CDATA[Armure]]></class><subclass id="2"><![CDATA[Armures en cuir]]></subclass><icon displayId="40493">inv_bracer_15</icon><inventorySlot id="9">Poignets</inventorySlot><htmlTooltip><![CDATA[<table><tr><td><table style="display:inline-table; vertical-align:inherit"><tr><td><!--nstart--><b class="q4">Brassards de malignité</b><!--nend--></td><th><b class="q0 whtt-extra">Phase 1</b></th></tr></table><!--ndstart--><!--ndend--><span class="q"><br>Niveau d'objet <!--ilvl-->115</span><!--bo--><br>Lié quand ramassé<!--ue--><table width="100%"><tr><td>Poignets</td><th><!--scstart4:2--><span class="q1">Cuir</span><!--scend--></th></tr></table><span><!--amr-->Armure : 159</span><br><span><!--stat7-->+25 Endurance</span><!--ebstats--><!--egstats--><!--eistats--><!--e--><!--ps--><br>Durabilité 40 / 40</td></tr></table><table><tr><td>Niveau <!--rlvl-->70 requis<br><span class="q2">Équipé\u00a0: Augmente de <!--rtg32-->22 le score de coup critique.</span><br><span class="q2">Équipé : <a href="https://fr.tbc.wowhead.com/spell=14056" class="q2">Augmente de 50 la puissance d'attaque.</a></span><!--itemEffects:1--><div class="whtt-sellprice">Prix de Vente: <span class="moneygold">2</span> <span class="moneysilver">80</span> <span class="moneycopper">84</span></div><div class="whtt-extra whtt-droppedby">Dépouillé sur: Damoiselle de vertu</div><div class="whtt-extra whtt-dropchance">Chance de Butin: 16.73%</div></td></tr></table><!--i?28514:1:70:70-->]]></htmlTooltip><json><![CDATA["appearances":{"0":[40493,""]},"armor":159,"classs":4,"displayid":40493,"flags2":24580,"id":28514,"level":115,"name":"Brassards de malignité","quality":4,"reqlevel":70,"slot":9,"slotbak":9,"source":[2],"sourcemore":[{"bd":1,"dd":1,"n":"Damoiselle de vertu","t":1,"ti":16457,"z":3457}],"subclass":2]]></json><jsonEquip><![CDATA["appearances":{"0":[40493,""]},"armor":159,"displayid":40493,"dura":40,"mleatkpwr":50,"mlecritstrkrtng":22,"reqlevel":70,"rgdatkpwr":50,"rgdcritstrkrtng":22,"sellprice":28084,"slotbak":9,"sta":25]]></jsonEquip><link>https://fr.tbc.wowhead.com/item=28514</link></item></wowhead>`))
	item, err := wowhead.GetItem(context.TODO(), 28514)
	td.CmpNoError(t, err)
	td.Cmp(t, item.Name, "Brassards de malignité")
	td.Cmp(t, item.Slot, "Poignets")
	td.Cmp(t, item.Sockets, int64(0))
}
