package initall

import (
	"fmt"
	"testing"
)

var (
	errOlives = fmt.Errorf("olives failed")
	errFigs   = fmt.Errorf("figs failed")
	errWine   = fmt.Errorf("wine failed")
)

type olives struct {
	color string
}

func (o *olives) InitializeOnce() error {
	return errOlives
}

type figs struct {
	color string
}

func (o *figs) InitializeOnce() error {
	return errFigs
}

type wine struct {
	color string
}

func (o *wine) InitializeOnce() error {
	return nil
}

func TestAll(t *testing.T) {
	initOlives := NewRegistrant()
	initOlives.Register(&olives{color: "green"})

	initFigs := NewRegistrant()
	initFigs.Register(&figs{color: "black"})

	initWine := NewRegistrant()
	initWine.Register(&wine{color: "red"})

	initMain := NewInitAll()

	errs := initMain.RunAllOnce()
	if errs == nil {
		t.Fatalf(`RunAllOnce didnt receive error`)
	}

	if initMain.HasError(errOlives) == nil {
		t.Fatalf(`RunAllOnce expected error for olives: received error: %v`, errOlives)
	}

	if initMain.HasError(errFigs) == nil {
		t.Fatalf(`RunAllOnce expected error for figs: received error: %v`, errFigs)
	}

	// should not be an error for wine.
	if initMain.HasError(errWine) != nil {
		t.Fatalf(`RunAllOnce didnt expect error for wine: received error`)
	}

	errs = initMain.Errors()
	if errs == nil {
		t.Fatalf(`Error didnt receive error`)
	}
}
