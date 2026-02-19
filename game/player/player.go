package player

import (
	"fmt"
	. "github.com/oms125/epoch-bot/game/items"
	. "github.com/oms125/epoch-bot/game"
)

type Player struct {
	ID string
	Lvl int
	InvSize int
	ArmSize int

	Inv Inventory
	Arm Armory
	Equipment
}

func (p *Player) GenerateAttack() *Attack {
	return &Attack{}
}

func (p *Player) ProcessAttack() {
	
}

func (p *Player) AddItem(id int, quantity ...int) error {
	data := ITEM_DATA[id]
	switch data.Type() {
	case TYPE_MATERIAL:
		num := 1
		if len(quantity) > 0 {
			num = quantity[0]
		}
		return p.AddMaterial(id, num)
	case TYPE_GEAR:
		return p.AddGear(ITEM_DATA[id].(*GearData).NewGear())
	}
	return fmt.Errorf("Invalid item ID")
}

func (p *Player) AddMaterial(id int, quantity int) error {
	return p.Inv.AddMaterial(id, quantity, p.InvSize)
}

func (p *Player) RemoveMaterial(id int, quantity int) error {
	return p.Inv.RemoveMaterial(id, quantity)
}

func (p *Player) AddGear(g *Gear) error {
	return p.Arm.AddGear(g, p.ArmSize)
}

func (p *Player) RemoveGear(idx int) {
	p.Arm.RemoveGear(idx)
}

func (p *Player) Equip(armSlot int, slot string) error {
	if (armSlot < len(p.Arm)) {
		g := p.Arm[armSlot]
		return p.Equipment.Equip(g, slot)
	}
	return fmt.Errorf("There is not item in armory slot **%d**", armSlot)
}

func (p *Player) Unequip(slot string) {
	p.Equipment.Unequip(slot)
}

func (p *Player) BuildEquipment() {
	equipment := make(map[string]*Gear)
	for _, t := range p.Arm {
		if t.Equipped != "" {
			equipment[t.Equipped] = t
		}
	}
	p.Equipment = equipment
}
