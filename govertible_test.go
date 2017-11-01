package govertible

import (
	"encoding/json"
	"testing"
)

type TypeASubset struct {
	SubSetStr    string
	SubSrtStrPtr *string
}

func newTypeASubset() *TypeASubset {
	StrPtr := "str"

	return &TypeASubset{
		SubSetStr:    "str",
		SubSrtStrPtr: &StrPtr,
	}
}

type TypeA struct {
	TypeASubset

	Str    string
	I      int
	Float  float32
	Double float64

	StrPtr    *string
	IPtr      *int
	FloatPtr  *float32
	DoublePtr *float64

	ArrStrs       []string
	ArrStrPtrs    []*string
	ArrPtrStrs    *[]string
	ArrPtrStrPtrs *[]*string

	MapStrStr map[string]string
}

func newTypeA() *TypeA {

	strPtr := "str"
	iPtr := 1
	floatPtr := float32(100.0)
	doublePtr := float64(100.0)
	arrPtrStrs := []string{"str"}
	arrPtrStrPtrs := []*string{&strPtr}

	return &TypeA{
		TypeASubset:   *newTypeASubset(),
		Str:           "str",
		I:             100,
		Float:         100.0,
		Double:        100.0,
		StrPtr:        &strPtr,
		IPtr:          &iPtr,
		FloatPtr:      &floatPtr,
		DoublePtr:     &doublePtr,
		ArrStrs:       []string{"str"},
		ArrStrPtrs:    []*string{&strPtr},
		ArrPtrStrs:    &arrPtrStrs,
		ArrPtrStrPtrs: &arrPtrStrPtrs,
	}
}

type TypeAClone struct {
	TypeA
}

type TypeASuperset struct {
	TypeA

	ssStr    string
	ssStrPtr *string
}

func newTypeASuperset() *TypeASuperset {

	ssStrPtr := "str"

	return &TypeASuperset{
		ssStr:    "str",
		ssStrPtr: &ssStrPtr,
	}
}

// check pointer mismatch
type TypeB struct {
	Str string
}

func newTypeB() *TypeB {
	return &TypeB{
		Str: "str",
	}
}

type TypeC struct {
	Str *string
}

func newTypeC() *TypeC {

	strPtr := "str"

	return &TypeC{
		Str: &strPtr,
	}
}

// check object mismatch
type TypeD struct {
	Obj1 TypeA
	Obj2 TypeB
}

func newTypeD() *TypeD {
	return &TypeD{
		Obj1: *newTypeA(),
		Obj2: *newTypeB(),
	}
}

type TypeE struct {
	Obj1 TypeB
	Obj2 *TypeA
}

func newTypeE() *TypeE {
	return &TypeE{
		Obj1: *newTypeB(),
		Obj2: newTypeA(),
	}
}

func TestClone(test *testing.T) {
	typeA := newTypeA()
	typeAClone := &TypeAClone{}

	if err := Convert(typeA, typeAClone); err != nil {
		test.Error("failed to convert typeA -> typeAClone: " + err.Error())
	}
}

func TestSuperset(test *testing.T) {
	typeA := newTypeA()
	typeASuperset := newTypeASuperset()

	if err := Convert(typeA, typeASuperset); err != nil {
		test.Error("failed to convert typeA -> typeASuperset: " + err.Error())
	}

	//test.Log(marshalJSON(typeASuperset))
}

func TestSubset(test *testing.T) {
	typeA := newTypeA()
	typeASubset := newTypeASubset()

	if err := Convert(typeA, typeASubset); err != nil {
		test.Error("failed to convert typeA -> typeASubset: " + err.Error())
	}

	//test.Log(marshalJSON(typeASubset))
}

func TestValToPtrConverions(test *testing.T) {
	typeB := newTypeB()
	typeC := &TypeC{}

	if err := Convert(typeB, typeC); err != nil {
		test.Error("failed to convert typeB -> typeC: " + err.Error())
	}

	if *typeC.Str != "str" {
		test.Error("typeC Str field not set")
	}

	//test.Log(marshalJSON(typeA))
	//test.Log(marshalJSON(TypeAPtrConversions))
}

func TestPtrToValConverions(test *testing.T) {
	typeC := newTypeC()
	typeB := &TypeB{}

	if err := Convert(typeC, typeB); err != nil {
		test.Error("failed to convert typeA -> typeAClone: " + err.Error())
	}

	if typeB.Str != "str" {
		test.Error("typeB Str field not set")
	}

	//test.Log(marshalJSON(typeA))
	//test.Log(marshalJSON(TypeAPtrConversions))
}

func TestStructFieldConverions(test *testing.T) {
	typeD := newTypeD()
	typeE := &TypeE{}

	if err := Convert(typeD, typeE); err != nil {
		test.Error("failed to convert typeD -> typeD: " + err.Error())
	}

	//test.Log(marshalJSON(typeD))
	//test.Log(marshalJSON(typeE))
}

func TestStructArrayPtr(test *testing.T) {
	typeBs := []*TypeB{newTypeB(), newTypeB()}
	typeCs := []*TypeC{}

	if err := Convert(typeBs, &typeCs); err != nil {
		test.Error("failed to convert []*typeB -> []*typeC: " + err.Error())
	}

	//test.Log(marshalJSON(typeBs))
	//test.Log(marshalJSON(typeCs))
}

func TestStructArray(test *testing.T) {
	typeBs := []TypeB{*newTypeB(), *newTypeB()}
	typeCs := []TypeC{}

	if err := Convert(typeBs, &typeCs); err != nil {
		test.Error("failed to convert []typeB -> []typeC: " + err.Error())
	}

	//test.Log(marshalJSON(typeBs))
	//test.Log(marshalJSON(typeCs))
}

func BenchmarkClone(b *testing.B) {
	for n := 0; n < b.N; n++ {
		typeA := newTypeA()
		typeAClone := &TypeAClone{}
		Convert(typeA, typeAClone)
	}
}

func marshalJSON(val interface{}) string {
	bytes, _ := json.MarshalIndent(val, "", "    ")
	return string(bytes)
}
