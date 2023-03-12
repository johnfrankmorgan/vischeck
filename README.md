# vischeck

A simple linter which checks certain visibility rules have been adhered to.


## Usage

Visibility rules are defined as tags on struct fields. For example:

```go
package main

type MyStruct struct {
	Name string `visibility:"readonly"`
	Age  int    `visibility:"readonly"`
}

func main() {
	s := MyStruct{}

	s.Name = "Frank"  // misuse of readonly field: cannot assign
	_ = &s.Name       // misuse of readonly field: cannot take address
	s.Age++           // misuse of readonly field: cannot increment
}
```


### Available Rules

The available rules are as follows:


#### `readonly`

Checks fields marked as `readonly` are only mutated by receiver functions
defined on their containing type.


#### Limitations

* Fields stored as a pointer type cannot be checked, instead `vischeck` will
  emit an error indicating that the field type is invalid.
* Struct literals are not yet checked.


## Rationale

In order to simplify (un)marshalling structs, it's often convenient to
export fields. Typically it is safe to read these fields but setting them may
require using a specific API in order to prevent invalid states. Ensuring the
correct APIs are used to set these fields is usually the responsibility of the
developer and code reviewer - this is where `vischeck` comes in. With
`vischeck`, developers will get a linting error if they attempt to mutate
certain struct fields directly, forcing them to refactor their code to use the
correct API. This linter would be disabled where direct write access to fields
is required (for example in some data access layer).

```go
package models

import (
	"errors"
	"fmt"
)

type Order struct {
	// The visibility tag ensures this field is never manipulated directly,
	// forcing developers to use the `SetStatus` method, instead.
	Status OrderStatus `visibility:"readonly"`
}

var ErrInvalidOrderStatus = errors.New("invalid status")

func (order *Order) SetStatus(status OrderStatus) error {
	if !status.Valid() {
		return fmt.Errorf("%w: %s", ErrInvalidOrderStatus, status)
	}

	if status <= order.Status || status != order.Status+1 {
		return fmt.Errorf("%w: can't transition from %s to %s", ErrInvalidOrderStatus, order.Status, status)
	}

	order.Status = status
	return nil
}

type OrderStatus int

const (
	OrderStatusInvalid OrderStatus = iota
	OrderStatusPending
	OrderStatusPaid
	OrderStatusShipped
	OrderStatusComplete
)

func (status OrderStatus) Valid() bool {
	return status > OrderStatusInvalid && status <= OrderStatusComplete
}

func (status OrderStatus) String() string {
	return "" // imagine some code here
}
```


## Building

To build `vischeck`, simply run `make build`. This will produce a standalone
binary (`build/vischeck`) and a plugin which can be used with
[golangci-lint](https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint).
(`build/vischeck.so`).
