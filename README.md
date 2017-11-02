# govertible

![logo](https://github.com/obsiius/govertible/logo.png)

[![godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/obsius/govertible)
[![coverage](https://coveralls.io/repos/github/obsius/govertible/badge.svg?branch=master)](https://coveralls.io/github/obsius/govertible?branch=master)
[![go report](https://goreportcard.com/badge/obsius/govertible)](https://goreportcard.com/report/obsius/govertible)
[![build](https://travis-ci.org/obsius/govertible.svg?branch=master)](https://travis-ci.org/obsius/govertible)
[![license](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/obsius/govertible/master/LICENSE)

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
## Examples
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
		Name:  "El Capitan",
		Phone: []byte("123-456-7890"),
		Age:   20,
	}
	employee := employee{}

	// convert a person struct to an employee struct
	govertible.Convert(&person, &employee)

	fmt.Println(employee)
}
```

Advanced example converting one struct to another using the convertTo interface.
```golang
package main

import (
	"fmt"
	"github.com/obsius/govertible"
)

type employee struct {
	Name  *string
	Alias string
	Phone []byte
	ID    uint64
}
type person struct {
	Name  string
	Alias string
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
		Name:  "El Capitan",
		Alias: "The Chief",
		Phone: []byte("123-456-7890"),
		Age:   10,
	}
	employee := employee{}

	// convert a person struct to an employee struct
	govertible.Convert(&person, &employee)

	fmt.Println(employee)
}
```

Advanced example converting one struct to another using the convertFrom interface.
```golang
package main

import (
	"fmt"
	"github.com/obsius/govertible"
)

type employee struct {
	Name  *string
	Alias string
	Phone []byte
	ID    uint64
}
type person struct {
	Name  string
	Alias string
	Phone []byte
	Age   int
}

func (this *employee) ConvertFrom(val interface{}) (bool, error) {
	switch val.(type) {
	case *person:
		v := val.(*person)
		govertible.ConvertFields(val, this)
		this.ID = uint64(v.Age)
		break
	}

	return false, nil
}

func main() {
	person := person{
		Name:  "El Capitan",
		Alias: "The Chief",
		Phone: []byte("123-456-7890"),
		Age:   10,
	}
	employee := employee{}

	// convert a person struct to an employee struct
	govertible.Convert(&person, &employee)

	fmt.Println(employee)
}
```

## Benchmarks
operation|ns/op|# operations|total time
-|-|-|-
ConvertStruct|22,850|100k|2.5s