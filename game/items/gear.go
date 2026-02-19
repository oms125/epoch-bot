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
		ImageName 	  string
		MaxDurability int
		EquipSlot     string
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

func (g *Gear) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, ":%s: **%s**\n", g.ImageName, g.Name)
	if (g.Attack > 0) {
		fmt.Fprintf(&str, "Attach +%d\n", g.Attack)
	}
	if (g.Defense > 0) {
		fmt.Fprintf(&str, "Defense +%d\n", g.Defense)
	}
	if (g.Health > 0) {
		fmt.Fprintf(&str, "Health +%d\n", g.Health)
	}
	return str.String()
}

//Armory
func (a Armory) AddGear(gear *Gear, armSize int) error {
	if len(a) >= armSize {
		return &ArmFullError{}
	}
	a = append(a, gear)
	return nil
}

func (a Armory) RemoveGear(idx int) {
	if (idx < len(a)) {
		a = slices.Delete(a, idx, idx+1)
	}
}

func (arm Armory) String() string {
	var armString strings.Builder; 
	for _, tool := range arm {
		fmt.Fprintf(&armString, "%s: (%d/%d) %s\n", tool.Name, tool.Durability, tool.MaxDurability, tool.Equipped)
	}
	return armString.String()
}

//Equipment
func (e Equipment) Equip(g *Gear, slot string) error {
	slot = strings.ToLower(slot)
	switch g.GearData.EquipSlot {
	case slot: e.equip(g, slot)
	case "tool":
		if (slot == "primary" || slot == "secondary") {
			e.equip(g, slot)
		}
	default:
		return fmt.Errorf("Invalid equipment slot for **%s**", g.Name)
	}
	return nil
}

func (e Equipment) equip(g *Gear, slot string) {
	e[slot].Equipped = ""
	e[slot] = g
	g.Equipped = slot
}

func (e Equipment) Unequip(slot string) {
	e[slot].Equipped = ""
	e[slot] = nil
}