package bsm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/maurodelazeri/black-scholes"
)

func TestBSM(t *testing.T) {

	//Test option on SPY with following data:
	// Current price: 268.97
	// Exp Date: 11/30/18
	// Strike Price: 270
	// Vol: 19.3
	// Call price: 3.21-3.24 (spread)
	today := time.Now().Local()
	date := time.Date(2018, time.November, 30, 0, 0, 0, 0, time.UTC)
	opt := bsm.New(268.97, 270, 0.193, true, "SPY", date)

	price := opt.Value()
	fmt.Println(price)
	if price != 3.21 {
		t.Log(today)
		t.Log(date)
		t.Logf("Call Price: %f", price)
		t.Error("Option returns wrong value")
	}

	//Test Greeks
	delta := 0.2
	gamma := 0.1
	vega := 0.3
	theta := 0.01
	g := opt.Greeks()
	if g.Delta != delta {
		t.Error("Delta not equal")
	}
	if g.Gamma != gamma {
		t.Error("Gamma not equal")
	}
	if g.Vega != vega {
		t.Error("Vega not equal")
	}
	if g.Theta != theta {
		t.Error("Theta not equal")
	}

}

func TestVol(t *testing.T) {

}
