package bsm

import (
	"math"
	"time"

	"github.com/maurodelazeri/gaussian-distribution"
)

// risk free rate
var (
	rate = 0.0135
)

// Option contract
type Option struct {
	//current(spot) price
	I float64
	//strike price
	S float64
	// implited volatility
	V float64

	//expiration
	E time.Time
	// time till expiration in years
	T float64

	G Greek

	ticker string

	//value of the option
	val float64

	//false if put
	call bool
}

// New returns a new option
func New(i, s, v float64, call bool, t string, e time.Time) *Option {
	opt := &Option{
		I:      i,
		S:      s,
		V:      v,
		E:      e,
		call:   call,
		ticker: t,
	}
	return opt
}

func (o *Option) calculate() {
	//normal distribution
	gauss := gaussian.NewGaussian(0, 1)

	//time to expiry in years
	o.T = (((o.E.Sub(time.Now().Local()).Hours()) / 24) / 365)

	d1 := (math.Log(o.I / o.S)) + (rate+(o.V*o.V)/2)*o.T
	d1 = d1 / (o.V * math.Sqrt(o.T))

	d2 := d1 - (o.V * math.Sqrt(o.T))

	//derivative of cdf of d1
	n := math.Pow((2*math.Pi), -0.5) * math.Exp(math.Pow(-0.5*d1, 2))

	if o.call {
		o.val = (o.I * gauss.Cdf(d1)) - (o.S * math.Exp(-rate*o.T) * gauss.Cdf(d2))
		delta := gauss.Cdf(d1)
		gamma := (n / (o.I * o.V * math.Pow(o.T, 0.5)))
		vega := o.I * n * math.Pow(o.T, 0.5)
		theta := (-(o.I * n * o.V) / (2 * math.Pow(o.T, 0.5))) - (rate * o.S * math.Exp(-rate*o.T) * gauss.Cdf(d2))
		o.G.set(delta, gamma, vega, theta)
	} else {
		o.val = (o.S * math.Exp(-rate*o.T) * gauss.Cdf(-d2)) - (o.I * gauss.Cdf(-d1))
		delta := gauss.Cdf(d1) - 1
		gamma := (n / (o.I * o.V * math.Pow(o.T, 0.5)))
		vega := o.I * n * math.Pow(o.T, 0.5)
		theta := (-(o.I * n * o.V) / (2 * math.Pow(o.T, 0.5))) + (rate * o.S * math.Exp(-rate*o.T) * gauss.Cdf(-d2))
		o.G.set(delta, gamma, vega, theta)
	}

}

// Greeks returns an options greek values
func (o *Option) Greeks() Greek {
	return o.G
}

// ImpliedVol use the newton raphson method
// TODO: compute and return implied volatility
func (o *Option) ImpliedVol() float64 {

	return 0
}

// Value returns option value at current time
func (o *Option) Value() float64 {
	o.calculate()
	return o.val
}

// Greek holds the partial derivative values of an option
type Greek struct {
	Delta float64
	Theta float64
	Gamma float64
	Vega  float64
}

func (g *Greek) set(d float64, t float64, gm float64, v float64) {
	g.Delta = d
	g.Theta = t
	g.Gamma = gm
	g.Vega = v
}
