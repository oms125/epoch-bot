package items

import (

)

//Stat calculations
func (g *Gear) HealthStat() int {
	return g.Health
}

func (g *Gear) DefenseStat() int {
	return g.Defense
}

func (g *Gear) AttackStat() int {
	return g.Attack
}