package assert

// Copygen defines the functions that will be generated.
type Copygen interface {
	AssertHuman(Human) Animal
	AssertPointer(*Human) Animal
	AssertInterface(Animal) Human
	AssertInterfacePointer(*Animal) Human
}
