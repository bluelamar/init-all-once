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
	// singleInstance = nil

	AddRegistrant(&olives{color: "black"}, 10) // Middle priority

	AddRegistrant(&figs{color: "green"}, 5) // Lowest priority

	AddRegistrant(&wine{color: "rose"}, 20) // Highest priority

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
		t.Fatalf(`Errors didnt receive error`)
	}

	if len(errs) != 2 {
		t.Fatalf("Expected 2 errors, but got %d errors", len(errs))
	}

	if errs[0] != errOlives {
		t.Fatalf("Expected highest priority olives: want olives error but got %v", errs[0])
	}

	if errs[1] != errFigs {
		t.Fatalf("Expected lowest priority figs: want figs error but got %v", errs[1])
	}
}
