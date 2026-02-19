package game

type Entity interface {
	GenerateAttack() *Attack
	ProcessAttack()
}

type Attack struct {

}