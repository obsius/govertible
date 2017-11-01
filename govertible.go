package govertible

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	invalidToInterface = "Can't convert to non-addressable struct (pass a pointer)"
	failedToConvert    = "Failed to convert"
)

// ConvertableTo interface for convertibles.
type ConvertableTo interface {
	ConvertTo(interface{}) (bool, error)
}

// ConvertableFrom interface for convertibles.
type ConvertableFrom interface {
	ConvertFrom(interface{}) (bool, error)
}

// MustConvertFields copies all of the fields inside of the "from" interface into the "to" interface.
// Panics if the operation fails.
func MustConvertFields(from interface{}, to interface{}) {
	err := ConvertFields(from, to)
	if err != nil {
		panic(err)
	}
}

// ConvertFields copies all of the fields inside of the "from" interface into the "to" interface.
// Returns an error if the operation fails.
func ConvertFields(from interface{}, to interface{}) error {

	if err := enforcePointer(to); err != nil {
		return err
	}

	fromVal := getVal(from)
	toVal := getVal(to)

	for i := 0; i < fromVal.NumField(); i++ {

		fromFieldName := fromVal.Type().Field(i).Name

		fromField := fromVal.Field(i)
		toField := toVal.FieldByName(fromFieldName)

		// ignore if fields don't map
		if !toField.IsValid() {
			continue
		}

		// ignore if either can't address or get interface
		if !toField.CanAddr() || !toField.CanInterface() {
			continue
		}

		// deref pointers
		fromVal := getReflectVal(fromField)
		toVal := getReflectVal(toField)

		fromInterface := fromField.Interface()
		toInterface := toField.Interface()

		// don't convert if a ptr is nil
		if fromVal.Kind() == reflect.Invalid {
			continue
		}

		// initialize empty pointers
		if toVal.Kind() == reflect.Invalid {
			toField.Set(reflect.New(toField.Type().Elem()))
			toVal = toField.Elem()
			toInterface = toField.Elem().Addr().Interface()
		}

		// convert, or set, or attempt to manipulate pointers to work
		if fromVal.Kind() == reflect.Struct && toVal.Kind() == reflect.Struct {
			if err := Convert(fromInterface, toVal.Addr().Interface()); err != nil {
				return err
			}
		} else if fromVal.Kind() == toVal.Kind() {
			if reflect.TypeOf(fromInterface) == reflect.TypeOf(toInterface) {
				toField.Set(fromField)
			} else if fromField.Kind() == reflect.Ptr {
				toField.Set(fromVal)
			} else if toField.Kind() == reflect.Ptr {
				toField.Set(fromField.Addr())
			}
		} else {
			fmt.Println("ignoring")
		}
	}

	return nil
}

// MustConvert copies values within the "from" interface into the "to" interface.
// Panics if the operation fails.
func MustConvert(from interface{}, to interface{}) {
	err := Convert(from, to)
	if err != nil {
		panic(err)
	}
}

// Convert copies values within the "from" interface into the "to" interface
// Returns an error if the operation fails.
func Convert(from interface{}, to interface{}) error {

	if err := enforcePointer(to); err != nil {
		return err
	}

	var err error

	toVal := reflect.ValueOf(to)

	fromElem := getVal(from)
	toElem := getVal(to)

	// deref all pointers
	for fromElem.Kind() == reflect.Ptr {
		fromElem = fromElem.Elem()
	}
	for toElem.Kind() == reflect.Ptr {
		toElem = toElem.Elem()
	}

	fromInterface := fromElem.Interface()
	toInterface := toElem.Addr().Interface()

	converted := false

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

	// if both from and to are structs, then convert fields
	if !converted && fromElem.Kind() == reflect.Struct && toElem.Kind() == reflect.Struct {
		if err := ConvertFields(from, to); err != nil {
			return err
		}
		converted = true
	}

	if !converted && fromElem.Kind() == reflect.Slice && toElem.Kind() == reflect.Slice {
		if err := convertArray(fromElem, toVal); err != nil {
			return err
		}
		converted = true
	}

	if !converted {
		return errors.New(failedToConvert)
	}

	return nil
}

func convertArray(fromVal reflect.Value, toVal reflect.Value) error {

	if !toVal.IsValid() {
		return nil
	}

	toValRef := toVal.Elem()

	if !toValRef.CanSet() {
		return nil
	}

	fromInterface := fromVal.Interface()
	toInterface := toValRef.Addr().Interface()

	fromKind := getType(fromInterface).Kind()
	toKind := getType(toInterface).Kind()

	if toKind == reflect.Array {
		fmt.Println("array")
	} else if toKind == reflect.Struct {
		for i := 0; i < fromVal.Len(); i++ {

			newToVal := reflect.New(getType(toInterface))
			newToInterface := newToVal.Interface()

			var fromInterface interface{}
			if fromVal.Index(i).Type().Kind() == reflect.Ptr {
				fromInterface = fromVal.Index(i).Interface()
			} else {
				fromInterface = fromVal.Index(i).Addr().Interface()
			}

			Convert(fromInterface, newToInterface)

			if getSliceType(toInterface).Kind() == reflect.Ptr {
				toVal.Elem().Set(reflect.Append(toVal.Elem(), newToVal))
			} else {
				toVal.Elem().Set(reflect.Append(toVal.Elem(), newToVal.Elem()))
			}
		}
	} else if toKind == fromKind {
		for i := 0; i < fromVal.Len(); i++ {
			fromInterface := fromVal.Index(i).Interface()
			toVal.Elem().Set(reflect.Append(toValRef, reflect.ValueOf(fromInterface)))
		}
	}

	return nil
}

func enforcePointer(to interface{}) error {
	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		return errors.New(invalidToInterface)
	}
	return nil
}

func getVal(val interface{}) reflect.Value {
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		return reflect.ValueOf(val).Elem()
	}
	return reflect.ValueOf(val)
}

func getReflectVal(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		return val.Elem()
	}
	return val
}

func getSliceType(a interface{}) reflect.Type {
	for t := reflect.TypeOf(a).Elem(); ; {
		switch t.Kind() {
		case reflect.Slice:
			t = t.Elem()
		default:
			return t
		}
	}
}

func getType(a interface{}) reflect.Type {
	for t := reflect.TypeOf(a); ; {
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice:
			t = t.Elem()
		default:
			return t
		}
	}
}
