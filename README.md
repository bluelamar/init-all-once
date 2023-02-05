# Priority Order Initialization of Multiple Components

priority-order-init-go

Allows multiple components of an application to be initialized in a user specified order.
Components with higher priority are initialized before components of lower priority.
The user determines the priority for each component to be initialized.

Golang has a limitation that a *variable* of type sync.once can only be run once in an application
despite being called with different functions. e.g.:

```go
var once sync.Once

once.Do(my_PackageA_Func1) 
once.Do(my_PackageA_Func2) 
once.Do(my_PackageB_Func) 
```

In the above example, only `my_PackageA_Func1` would be called.
Neither `my_PackageA_Func2` nor `my_PackageB_Func` would get called.

This package will allow an application to register multiple components
such that each component will be able to have initialization functionality
run once.
It allows the user to specify an ordering via a priority per initialization.
For component initializers that have the same priority, it is arbitrary
which initializer would run first.

## Usage

Each package that needs to initialize once can register itself with priority-order-init-go.
The registration can be done by each package leveraging their *init()* function(s).
Note that *init* functions are called only once per package.
The *init* functions run after all variable declarations and before the main function. 
Therefore each package will make use of their *init* functions for the registration process.

Components will implement the ComponentInitalizer interface for objects they use to
register initialization.

## Example

This example registers several objects from the same package.
It also registers another package in main() as a different use case.
Comments describe the priority of the initialization.

```go
package my_PackageA

import "github.com/bluelamar/priority-order-init-go/initall"

type MyA1 struct {
	// logger etc.
}

// InitializeOnce is the ComponentInitalizer interface implementation.
func (m *MyA1) InitializeOnce() error {
	// do something...
	return nil
}


type MyA2 struct {
	// logger etc.
}

// InitializeOnce is the ComponentInitalizer interface implementation.
func (m *MyA2) InitializeOnce() error {
	// do something...
	return nil
}

func init() {
	// Example where MyA2 will be initialized before MyA1 since it has higher priority.
	initall.AddRegistrant(&MyA1{...}, 10)
	initall.AddRegistrant(&MyA2{...}, 20)
}

...

package main

import (
	"my_PackageA" // my_PackageA registers itself via init()
	"my_PackageB" // my_PackageB does not register itself but does implement the ComponentInitalizer interface.

	"github.com/bluelamar/priority-order-init-go/initall"
)

func main() {
	...
	// Since my_PackageB does not register itself, here is an example of registering it from main.
	// It is of same priority as MyA1 from my_PackageA, so it is arbitrary whether it or MyA1 is initialized first.
	initall.AddRegistrant(my_PackageB.New(), 10)

	initAllOnce := initall.NewInitAll()

	if err := initAllOnce.RunAllOnce(); err != nil {
		log.Fatalf("init-all-once error: %s", err)
	}
    ...
}

```



