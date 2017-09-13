## govertible

A lightweight package to convert similar structures to and from each other.

govertible will convert matching fields names of identical types from a source to a destination struct.  `ConvertTo()` and `ConvertFrom()` are implementable and allow for custom conversions as shown in the examples below.

#### ConvertibleTo
```golang
type ConvertibleTo interface {
	ConvertTo(interface{}) (bool, error)
}
```

#### ConvertibleFrom
```golang
type ConvertibleTo interface {
	ConvertTo(interface{}) (bool, error)
}
```
## examples
Simple example converting one struct to another by field names without implementing an interface.
```golang
package main

import (
	"fmt"
	"github.com/obsius/govertible"
)

type employee struct {
	Name  string
	Phone []byte
	ID    uint64
}
type person struct {
	Name  string
	Phone []byte
	Age   int
}

func main() {
	person := person{
		Name:  "ElCaptain",
		Phone: []byte("123-456-7890"),
	}
	employee := employee{}

	// convert a person struct to an employee struct
	govertible.Convert(&person, &employee)

	fmt.Println(employee)
}
```

Advanced example converting one struct to another using an interface.
```golang
package main

import (
	"fmt"
	"github.com/obsius/govertible"
)

type employee struct {
	Name  string
	Phone []byte
	ID    uint64
}
type person struct {
	Name  string
	Phone []byte
	Age   int
}
func (this *person) ConvertTo(val interface{}) (bool, error) {
	switch val.(type) {
	case *employee:
		v := val.(*employee)
		govertible.ConvertFields(this, val)
		v.ID = uint64(this.Age)
		break
	}

	return false, nil
}

func main() {
	person := person{
		Name:  "ElCaptain",
		Phone: []byte("123-456-7890"),
		Age:   10,
	}
	employee := employee{}

	// convert a person struct to an employee struct
	govertible.Convert(&person, &employee)

	fmt.Println(employee)
}
```

## benchmarks
operation|ns/op|total time
-|-|-