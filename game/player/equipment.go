package player

import (
	"fmt"
	"strings"
)

type Equipment []*Tool

//Equipment Slots
const (
	_ = iota
	//Slots
	HEAD
	BODY
	LEGS
	FEET
	PRIMARY
	SECONDARY
	//Slot options only
	TOOL
	//Not a slot
	LENGTH
)

func (p *Player) BuildEquipment() {
	equipment := make([]*Tool, LENGTH)
	for _, t := range p.Arm {
		if t.Equipped > 0 {
			equipment[t.Equipped] = t
		}
	}
	p.Equipment = equipment
}

//Game Logic
func (p *Player) Equip(armorySlot int, equipmentSlot string) error {
	if armorySlot >= len(p.Arm) { return fmt.Errorf("There is nothing in armory slot **%d**", armorySlot) }
	tool := p.Arm[armorySlot]
	if tool == nil { return fmt.Errorf("There is nothing in armory slot **%d**", armorySlot) }
	slot := strings.ToLower(equipmentSlot)
	isSlotEmpty := slot == ""

	switch tool.ToolData.Slot {
	case HEAD: if isSlotEmpty || slot == "head" { p.equip(tool, HEAD); return nil }
	case BODY: if isSlotEmpty || slot == "body" { p.equip(tool, BODY); return nil }
	case LEGS: if isSlotEmpty || slot == "legs" { p.equip(tool, LEGS); return nil }
	case FEET: if isSlotEmpty || slot == "feet" { p.equip(tool, FEET); return nil }
	case TOOL: 
		if isSlotEmpty {
			return fmt.Errorf("Must specify equipment slot for **%s**", tool.Name)
		}
		if slot == "primary" { p.equip(tool, PRIMARY); return nil }
		if slot == "secondary" { p.equip(tool, SECONDARY); return nil }
	}
	return fmt.Errorf("Invalid equipment slot for **%s**", tool.Name)
}

func (p *Player) equip(t *Tool, slot int) {
	p.Unequip(t.Equipped)
	p.Unequip(slot)
	p.Equipment[slot] = t
	t.Equipped = slot
}

func (p *Player) Unequip(slot int) {
	t := p.Equipment[slot]
	if t == nil { return }
	t.Equipped = 0
	p.Equipment[slot] = nil

}