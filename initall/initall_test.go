// Copyright 2022, Initialize All Once Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
