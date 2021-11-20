package wowhead_test

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/origin-finkle/logs/internal/wowhead"
)

func TestGetGem(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(wowhead.Client)
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("GET", "https://fr.tbc.wowhead.com/item=1&xml", httpmock.NewStringResponder(200, `<?xml version="1.0" encoding="UTF-8"?><wowhead><error>Item not found!</error></wowhead>`))
	_, err := wowhead.GetGem(context.TODO(), 1)
	td.CmpError(t, err)

	httpmock.RegisterResponder("GET", "https://fr.tbc.wowhead.com/item=24030&xml", httpmock.NewStringResponder(200, `<?xml version="1.0" encoding="UTF-8"?><wowhead><item id="24030"><name><![CDATA[Rubis vivant runique]]></name><level>70</level><quality id="3">Rare</quality><class id="3"><![CDATA[Gemmes]]></class><subclass id="0"><![CDATA[Gemmes rouges]]></subclass><icon displayId="0">inv_jewelcrafting_livingruby_03</icon><inventorySlot id="0"></inventorySlot><htmlTooltip><![CDATA[<table><tr><td><table style="display:inline-table; vertical-align:inherit"><tr><td><!--nstart--><b class="q3">Rubis vivant runique</b><!--nend--></td><th><b class="q0 whtt-extra">Phase 1</b></th></tr></table><!--ndstart--><!--ndend--><span class="q whtt-extra whtt-ilvl"><br>Niveau d'objet <!--ilvl-->70</span><!--bo--><!--ue--><!--ebstats--><!--egstats--><!--eistats--><br /><span class="q1"><a href="https://fr.tbc.wowhead.com/spell=9415">+9 aux dégâts des sorts</a></span></td></tr></table><table><tr><td><span class="q">&quot;Correspond à une châsse rouge.&quot;</span><div class="whtt-sellprice">Prix de Vente: <span class="moneygold">3</span></div></td></tr></table>]]></htmlTooltip><json><![CDATA["classs":3,"flags2":24580,"id":24030,"level":70,"name":"Rubis vivant runique","quality":3,"slot":0,"source":[1],"sourcemore":[{"c":11,"icon":"inv_jewelcrafting_livingruby_03","n":"Rubis vivant runique","s":755,"t":6,"ti":31088}],"subclass":0]]></json><jsonEquip><![CDATA["avgbuyout":1254974,"sellprice":30000,"spldmg":9,"splheal":9]]></jsonEquip><createdBy><spell id="31088" name="Rubis vivant runique" icon="inv_jewelcrafting_livingruby_03" minCount="0" maxCount="0"><reagent id="23436" name="Rubis vivant" quality="3" icon="inv_jewelcrafting_livingruby_02" count="1" /></spell></createdBy><link>https://fr.tbc.wowhead.com/item=24030</link></item></wowhead>`))
	gem, err := wowhead.GetGem(context.TODO(), 24030)
	td.CmpNoError(t, err)
	td.Cmp(t, gem.Color, "red")
	td.Cmp(t, gem.Quality, int64(3))
	td.Cmp(t, gem.Name, "Rubis vivant runique")
}
