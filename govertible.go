package govertible

import (
	_ "../utility"

	"errors"
	_ "fmt"
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

func ConvertFields(from interface{}, to interface{}) error {

	if reflect.ValueOf(from).Kind() != reflect.Ptr {
		panic("Can't convert from non-addressable struct (pass a pointer)")
		return errors.New("Can't convert from non-addressable struct (pass a pointer)")
	}

	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		panic("Can't convert to non-addressable struct (pass a pointer)")
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

		//	fmt.Printf("IN FIELD: %v\n", fromFieldName)

		fromVal := rFrom.Field(i)
		toVal := rTo.FieldByName(fromFieldName)

		// ignore if fields don't map
		if !toVal.IsValid() {
			continue
		}

		// deref all pointers
		fromElem := fromVal
		toElem := toVal
		for fromElem.Kind() == reflect.Ptr {
			fromElem = fromElem.Elem()
		}
		for toElem.Kind() == reflect.Ptr {
			toElem = toElem.Elem()
		}

		//fmt.Printf("field: %v - %v - %v\n", fromFieldName, fromVal, toVal)

		fromInterface := fromVal.Addr().Interface()
		toInterface := toVal.Addr().Interface()

		//fmt.Printf("field: %v - %v - %v\n", fromFieldName, utility.MustMarshalJSON(fromInterface), utility.MustMarshalJSON(toInterface))

		// don't convert if a ptr is nil
		if fromElem.Kind() == reflect.Invalid {
			continue
		}

		// initialize empty pointers
		if toElem.Kind() == reflect.Invalid {
			toVal.Set(reflect.New(toVal.Type().Elem()))
			toElem = toVal.Elem()
			toInterface = toVal.Elem().Addr().Interface()
		}

		// same kind, then just set the val, else check for kinds
		if reflect.TypeOf(fromInterface) == reflect.TypeOf(toInterface) {
			toVal.Set(fromVal)
		} else {
			//	fmt.Printf("calling convert with  %v  %v\n", fromInterface, toInterface)
			Convert(fromInterface, toInterface)
		}
	}

	return nil
}

func Convert(from interface{}, to interface{}) error {

	var err error

	if reflect.ValueOf(from).Kind() != reflect.Ptr {
		panic("Can't convert from non-addressable struct (pass a pointer)")
		return errors.New("Can't convert from non-addressable struct (pass a pointer)")
	}

	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		panic("Can't convert to non-addressable struct (pass a pointer)")
		return errors.New("Can't convert to non-addressable struct (pass a pointer)")
	}

	//fromVal := reflect.ValueOf(from)
	toVal := reflect.ValueOf(to)

	fromElem := reflect.ValueOf(from).Elem()
	toElem := reflect.ValueOf(to).Elem()

	// deref all pointers
	for fromElem.Kind() == reflect.Ptr {
		fromElem = fromElem.Elem()
	}
	for toElem.Kind() == reflect.Ptr {
		toElem = toElem.Elem()
	}

	fromInterface := fromElem.Addr().Interface()
	toInterface := toElem.Addr().Interface()

	converted := false
	//fmt.Printf("IN CONVERT: TYPES: %v <----> %v\n", fromElem.Kind(), toElem.Kind())
	// if from is a struct, then try and call it's ConvertTo()
	if fromElem.Kind() == reflect.Struct {
		toConverter, ok := fromInterface.(ConvertableTo)
		if ok {
			converted, err = toConverter.ConvertTo(to)
			if err != nil {
				return err
			}
		}

	}
	// if unable, then try the to struct's ConvertFrom()
	if !converted && toElem.Kind() == reflect.Struct {
		fromConverter, ok := toInterface.(ConvertableFrom)
		if ok {
			converted, err = fromConverter.ConvertFrom(from)
			if err != nil {
				return err
			}
		}
	}

	if !converted && fromElem.Kind() == reflect.Struct && toElem.Kind() == reflect.Struct {
		ConvertFields(from, to)
	}

	// TODO check for byte[]
	/*
		if !converted && fromElem.Kind() == reflect.Slice && toElem.Kind() == reflect.Struct {
			fromType := getElemType(fromInterface)
			if fromType == uint8() {

			}
		}*/

	if !converted && fromElem.Kind() == reflect.Slice && toElem.Kind() == reflect.Slice {
		//toVal.Set(reflect.New(toVal.Type().Elem()))
		//toElem = toVal.Elem()
		//toInterface = toVal.Elem().Addr().Interface()

		//fmt.Println(reflect.New(getElemType(fromInterface)))

		s := fromElem
		for i := 0; i < s.Len(); i++ {
			//	fmt.Println("======>")

			newTo := reflect.New(getElemType(toInterface))

			newToInterface := newTo.Interface()

			thing := s.Index(i).Elem().Addr().Interface()

			Convert(thing, newToInterface)

			toVal.Elem().Set(reflect.Append(toElem, newTo))
		}

	}
	/*
		fmt.Println("--------------------------")
		fmt.Println("from")
		fmt.Println(utility.MustMarshalJSON(from))
		fmt.Println("to")
		fmt.Println(utility.MustMarshalJSON(to))
		fmt.Println("\n\n")
	*/
	return nil
}

func getElemType(a interface{}) reflect.Type {
	for t := reflect.TypeOf(a); ; {
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice:
			t = t.Elem()
		default:
			return t
		}
	}
}
