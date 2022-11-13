# Initialize All Once

Allows multiple components of an application to be initialized once.

Golang has a limitation that once.Do() will only be run once in an application
despite being called with different functions. e.g.:

```go
once.Do(my_A_PackageFunc) 
once.Do(my_B_PackageFunc) 
```

In the above example, only `my_A_PackageFunc` would be called.
`my_B_PackageFunc` would *not* get called.

This package will allow an application to register multiple components
such that each component will be able to have initialization functionality
run once.

## Usage

Each package that needs to initialize once must register itself with init-all-once-go.
This can be done by each package leveraging the init() function(s).
Init functions are called only once per package.
They run after all variable declarations and before the main function. 
Each package will make use of the init() function for the registration process.

Each package can register using their New constructor as well.

Example Registering from the package constructor function:

```go
package my_A_Package

import "github.com/bluelamar/init-all-once-go/initall"

type MyA struct {
	// logger etc.
}

// Many packages have a constructor function of some sort.
// This is just a simple example of one.
func New() *MyA {
	m := &MyA{...}

	// Registration STARTED.
	initAllOnceMyA = initall.NewRegistrant()

	// Registration COMPLETED.
	initAllOnceMyA.Register(m)

	return m
}

// InitializeOnce is the ComponentInitalizer interface implementation.
func (m *MyA) InitializeOnce() error {
	// do something...
	return nil
}
```

Example Registering from the package init() function:

```go
package my_B_Package

import "github.com/bluelamar/init-all-once-go/initall"

type myB_Init struct {
	// logger etc
}

// Package init function is called once for the package.
// Package must create a new registrant here.
func init() {
	my_B_PackageInit := &myB_Init{...}

	// Registration STARTED.
	initAllOnce_myB := initall.NewRegistrant()

	// Registration COMPLETED.
	initAllOnce_myB.Register(my_B_PackageInit)
}

// InitializeOnce is the ComponentInitalizer interface implementation.
func (m *myB_Init) InitializeOnce() error {
	// do something...
	return nil
}
```

Example of steps main takes:

```go
package main

import (
	"my_A_Package"
	"github.com/bluelamar/init-all-once-go/initall"
)

func main() {
	// my_B_Package has already registered via its init() function.

	initAllOnce := initall.NewInitAll()

	// my_A_Package registers here, before RunAllOnce() will be called.
	myA := my_A_Package.New()

	// RunAllOnce blocks until all registrants have called Register().
	if err := initAllOnce.RunAllOnce(); err != nil {
		log.Fatalf("init-all-once error: %s", err)
	}

	// ...
}

```

