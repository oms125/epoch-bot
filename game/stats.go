package game

type Stats interface {
	AttackStat()
	DefenseStat()
	HealthStat()
}

type BaseStats struct {
	Attack int
	Defense int
	Health int
}