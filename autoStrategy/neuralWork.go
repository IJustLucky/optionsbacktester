package autoStrategy

type Neuron struct {
	Activation bool
	Weight     []float64
	Delta      float64
	Bias       float64
}

type Layer struct {
	ActivationFunc func(float64) float64
	DerivativeFunc func(float64) float64
	Neurons        []*Neuron
}

type NeuralNetwork struct {
	Layers []*Layer
	Inputs int
}
