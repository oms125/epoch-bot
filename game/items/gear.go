package items

import (
	"fmt"
	"strings"
	"slices"
	. "github.com/oms125/epoch-bot/game"
)

const TYPE_GEAR = "gear"

type (
	Armory []*Gear
	Equipment map[string]*Gear

	Gear struct {
		*GearData
		GearMetadata
	}
	GearData struct {
		ID            int
		Name          string
		MaxDurability int
		EquipSlot     string
		Emoji
		BaseStats
	}
	GearMetadata struct {
		Durability int `db:"durability"`
		Equipped   string `db:"equipped"`
	}
)

var ArmoryFields = []Field{
	{
		Name: "item_id",
		Type: "INTEGER",
	},
	{
		Name: "durability",
		Type: "INTEGER",
	},
	{
		Name: "equipped",
		Type: "TEXT",
	},
}

//Gear
func (g *Gear) Type() string { return TYPE_GEAR }
func (g *GearData) Type() string { return TYPE_GEAR }

func (gd *GearData) NewGear() *Gear {
	return &Gear{
		GearData: gd,
		GearMetadata: GearMetadata{
			Durability: gd.MaxDurability,
			Equipped: "",
		},
	}
}

func (g *Gear) Info() string {
	var str strings.Builder
	fmt.Fprintf(&str, "%s **%s**\n", g.Emoji, g.Name)
	if (g.Attack > 0) {
		fmt.Fprintf(&str, "Attack +%d\n", g.Attack)
	}
	if (g.Defense > 0) {
		fmt.Fprintf(&str, "Defense +%d\n", g.Defense)
	}
	if (g.Health > 0) {
		fmt.Fprintf(&str, "Health +%d\n", g.Health)
	}
	return str.String()
}

func (g *Gear) String() string {
	return fmt.Sprintf("%s %s (%d/%d)", g.Emoji, g.Name, g.Durability, g.MaxDurability)
}

//Armory
func (a *Armory) AddGear(gear *Gear, armSize int) error {
	if len(*a) >= armSize {
		return &ArmFullError{}
	}
	*a = append(*a, gear)
	return nil
}

func (a *Armory) RemoveGear(idx int) {
	if (idx < len(*a)) {
		*a = slices.Delete(*a, idx, idx+1)
	}
}

func (a Armory) String() string {
	var armString strings.Builder; 
	for _, tool := range a {
		fmt.Fprint(&armString, tool.String(), "\n")
	}
	return armString.String()
}

//Equipment
func (e Equipment) Equip(g *Gear, slot string) error {
	slot = strings.ToLower(slot)
	switch g.GearData.EquipSlot {
	case slot: 
		e.equip(g, slot)
		return nil
	case "tool":
		if (slot == "primary" || slot == "secondary") {
			e.equip(g, slot)
			return nil
		}
	}
	return fmt.Errorf("Invalid equipment slot for **%s**", g.Name)
}

func (e Equipment) equip(g *Gear, slot string) {
	if (e[slot] != nil) {
		e[slot].Equipped = ""
	}
	e[slot] = g
	if (g.Equipped != "") {
		e[g.Equipped] = nil 
	}
	g.Equipped = slot
}

func (e Equipment) Unequip(slot string) {
	if (e[slot] != nil) {
		e[slot].Equipped = ""
		e[slot] = nil
	}
}

func (e Equipment) SlotValue(slot string) string {
	if (e[slot] != nil) {
		return e[slot].String()
	}
	return "Nothing Equipped"
}