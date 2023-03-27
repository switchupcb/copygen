package assert

// Animal represents an abstraction of an animal.
type Animal interface {
	Survive()
}

// Human represents an animal that belongs to the species Homo sapiens.
type Human struct{}

func (Human) Survive() {}
