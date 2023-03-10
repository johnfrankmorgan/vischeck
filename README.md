# Vis Check

A simple linter which checks certain visibility rules have been adhered to.

```go
package main

import "fmt"

type X struct {
	Y int `visibility:"readonly"`
}

func test(x X) {
	x.Y = 1000 // cannot assign
	x.Y++      // cannot increment
	_ = &x.Y   // cannot take address

	fmt.Println(x.Y) // this should be fine
}

func (x *X) SetY() {
	// all valid, in a receiver function
	x.Y = 2000
	x.Y++
	_ = &x.Y
}
```
