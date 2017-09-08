package govertible

import (
	"../utility"
	"errors"
	"fmt"
	"reflect"
)

// ConvertableTo interface for convertibles
type ConvertableTo interface {
	ConvertTo(interface{}) (bool, error)
}

// ConvertableFrom interface for convertibles
type ConvertableFrom interface {
	ConvertFrom(interface{}) (bool, error)
}

func Convert(from interface{}, to interface{}) error {

	var err error

	if !reflect.ValueOf(from).CanInterface() {
		return errors.New("Can't convert from non-addressable struct (pass a pointer)")
	}

	if !reflect.ValueOf(to).CanInterface() {
		return errors.New("Can't convert to non-addressable struct (pass a pointer)")
	}

	rFrom := reflect.ValueOf(from).Elem()
	rTo := reflect.ValueOf(to).Elem()

	// deref all pointers
	for rFrom.Kind() == reflect.Ptr {
		rFrom = rFrom.Elem()
	}
	for rTo.Kind() == reflect.Ptr {
		rTo = rTo.Elem()
	}

	for i := 0; i < rFrom.NumField(); i++ {

		fromFieldName := rFrom.Type().Field(i).Name

		fromVal := rFrom.Field(i)
		toVal := rTo.FieldByName(fromFieldName)

		// deref all pointers
		fromElem := fromVal
		toElem := toVal
		for fromElem.Kind() == reflect.Ptr {
			fromElem = fromElem.Elem()
		}
		for toElem.Kind() == reflect.Ptr {
			toElem = toElem.Elem()
		}

		fromInterface := fromVal.Addr().Interface()
		toInterface := toVal.Addr().Interface()

		fmt.Printf("field: %v - %v - %v\n", fromFieldName, utility.MustMarshalJSON(fromInterface), utility.MustMarshalJSON(toInterface))

		// don't convert if a ptr is nil
		if fromElem.Kind() == reflect.Invalid {
			continue
		}

		// initialize empty pointers
		if toElem.Kind() == reflect.Invalid {
			toVal.Set(reflect.New(toVal.Type().Elem()))
			toElem = toVal.Elem()
		}

		// same kind, then just set the val, else check for kinds
		if reflect.TypeOf(fromInterface) == reflect.TypeOf(toInterface) {
			toVal.Set(fromVal)
		} else {

			converted := false

			// if from is a struct, then try and call it's ConvertTo()
			if fromElem.Kind() == reflect.Struct {
				toConverter, ok := fromInterface.(ConvertableTo)
				if ok {
					converted, err = toConverter.ConvertTo(toInterface)
					if err != nil {
						return err
					}
				}

			}
			// if unable, then try the to struct's ConvertFrom()
			if !converted && toElem.Kind() == reflect.Struct {
				fromConverter, ok := toInterface.(ConvertableFrom)
				if ok {
					converted, err = fromConverter.ConvertFrom(fromInterface)
					if err != nil {
						return err
					}
				}
			}
			// if both failed, then just try a generic conversion
			if !converted && fromElem.Kind() == reflect.Struct && toElem.Kind() == reflect.Struct {
				Convert(fromInterface, toInterface)
			}
		}
	}

	fmt.Println("--------------------------")
	fmt.Println("from")
	fmt.Println(utility.MustMarshalJSON(from))
	fmt.Println("to")
	fmt.Println(utility.MustMarshalJSON(to))
	fmt.Println("\n\n")

	return nil
}
