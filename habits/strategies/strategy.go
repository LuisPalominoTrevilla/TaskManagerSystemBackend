package strategy

// Strategy represents an interface to the basic oprations of a model
type Strategy interface {
	CalculateScore(float32, int, int)
}