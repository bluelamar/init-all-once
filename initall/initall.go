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
	"sort"
	"sync"
)

// ComponentInitalizer is the interface that registering components should implement.
type ComponentInitalizer interface {
	InitializeOnce() error
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

type componentWithPriority struct {
	initializer ComponentInitalizer
	priority    int
}

// byCompPriority is used to sort the component initializers by priority.
type byCompPriority []*componentWithPriority

func (a byCompPriority) Len() int           { return len(a) }
func (a byCompPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCompPriority) Less(i, j int) bool { return a[i].priority > a[j].priority }

type initAllOnce struct {
	initCalls []*componentWithPriority
	errs      []error // Each component initializer could provide an error.
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

// AddRegistrant will add a component initializer with its relative priority.
// It should be called by the package/module init() function.
func AddRegistrant(compInit ComponentInitalizer, priority int) {
	newInitAll().register(compInit, priority)
}

func newInitAll() *initAllOnce {
	lock.Lock()
	defer lock.Unlock()

	if singleInstance == nil {
		singleInstance = &initAllOnce{
			initCalls: make([]*componentWithPriority, 0),
			errs:      make([]error, 0),
			doneOnce:  false,
		}
	}

	return singleInstance
}

func (i *initAllOnce) register(compInit ComponentInitalizer, priority int) {
	r := &componentWithPriority{
		initializer: compInit,
		priority:    priority,
	}

	i.initCalls = append(i.initCalls, r)
}

// RunAllOnce is called by the main function to run all registrants initializers.
func (i *initAllOnce) RunAllOnce() []error {
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
	// Sort the registrants by priority.
	sort.Sort(byCompPriority(i.initCalls))

	// Run the InitializeOnce() for each registrant in priority order.
	for _, e := range i.initCalls {
		if err := e.initializer.InitializeOnce(); err != nil {
			// log.Printf("initallonce:runall: elem returned %v\n", err)
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
