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
	"errors"
	"fmt"
	"sync"
)

// ComponentInitalizer is the interface that registering components should implement.
type ComponentInitalizer interface {
	InitializeOnce() error
}

// InitRegistrantOnce is the interface used by packages to register themselves.
type InitRegistrantOnce interface {
	// Register is called by each module to regiter themselves.
	Register(compInit ComponentInitalizer)
}

// InitAllOnce is the interface used by main().
type InitAllOnce interface {
	// RunAllOnce is called by main.
	RunAllOnce() []error

	// HasError returns the matched error if it exists.
	HasError(err error) error

	// Errors returns all errors set by RunAllOnce().
	Errors() []error
}

type initAllOnce struct {
	initCalls []ComponentInitalizer
	errs      []error // collect errors into one error
	doneOnce  bool
}

var (
	lock           = &sync.Mutex{}
	singleInstance *initAllOnce
	initWG         sync.WaitGroup
)

// NewInitAll is called by main(). Other packages should call NewRegistrant().
func NewInitAll() InitAllOnce {
	return newInitAll()
}

// NewRegistrant is called by each package except main.
// It should be called by the package/module init() function.
// Then each package/module must also call Register(compInit ComponentInitalizer).
func NewRegistrant() InitRegistrantOnce {
	initWG.Add(1)

	return newInitAll()
}

func newInitAll() *initAllOnce {
	lock.Lock()
	defer lock.Unlock()

	if singleInstance == nil {
		singleInstance = &initAllOnce{
			initCalls: make([]ComponentInitalizer, 0),
			errs:      make([]error, 0),
			doneOnce:  false,
		}
	}

	return singleInstance
}

// Register is called by registrants to specify their initializer.
func (i *initAllOnce) Register(compInit ComponentInitalizer) {
	i.initCalls = append(i.initCalls, compInit)
	initWG.Done()
}

// RunAllOnce is called by the main fucntion to run all registrants initializers.
func (i *initAllOnce) RunAllOnce() []error {
	// wait for registrars to finish.
	initWG.Wait()

	// Ensure we only run all initializers once for the application.
	lock.Lock()
	if i.doneOnce == true {
		lock.Unlock()
		return i.errs
	} else {
		i.doneOnce = true
	}
	lock.Unlock()

	i.runAll()

	return i.errs
}

func (i *initAllOnce) runAll() {
	// Run the InitializeOnce() for each registrant.
	for _, e := range i.initCalls {
		if err := e.InitializeOnce(); err != nil {
			fmt.Printf("FIX runall: elem returned %v\n", err)
			i.errs = append(i.errs, err)
		}
	}
}

func (i *initAllOnce) HasError(err error) error {
	for _, e := range i.errs {
		if errors.Is(e, err) {
			return e
		}
	}

	return nil
}

func (i *initAllOnce) Errors() []error {
	return i.errs
}
