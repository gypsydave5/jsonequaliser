package jsonequaliser

import (
	"fmt"
	"testing"
)

func TestItWorksWithEmptyArrays(t *testing.T) {
	A := `{"x": []}`
	B := `{"x": [1,2,3]}`

	if compat, err := IsCompatible(A, B); err != nil || !compat {
		t.Error("It broke, compat = ", compat, "err = ", err)
	}
}

func TestItAssumesFieldIsCompatibleIfAisNull(t *testing.T) {
	A := `{"max_score":null}`
	B := `{"max_score":10}`

	if compat, err := IsCompatible(A, B); err != nil || !compat {
		t.Error("It should've been compatible because we cant infer type from null = ", compat, "err = ", err)
	}
}

func TestItChecksKeyIsInBEvenWhenValueInAIsNull(t *testing.T) {
	A := `{"total":0,"max_score":null}`
	B := `{"total":0}`

	if compat, err := IsCompatible(A, B); err != nil || compat {
		t.Error("It should not be compatible because key is missing in B = ", compat, "err = ", err)
	}
}

func TestItWithSomethingRealish(t *testing.T) {
	A := `[{"MaterialID":"1234","ContentDate":{"From":"2005-04-20","To":"2015-12-01"},"AccessDate":{"From":"1987-04-20","To":"1990-12-01"}}]`
	B := `[{"MaterialID":"1234","ContentDate":{"From":"2005-04-20","To":"2015-12-01"},"AccessDate":{"From":"1987-04-20","To":"1990-12-01"}}]`

	if compat, err := IsCompatible(A, B); !compat {
		t.Log("Err", err)
		t.Error("NOT COMPATIBLE")
	}
}

func ExampleIsCompatible() {
	A := `[{"firstname": "chris", "lastname": "james", "age": 30}]`
	B := `[{"firstname": "Bob", "lastname": "Smith", "age": 25, "favourite-colour": "blue"}]`

	fmt.Println(IsCompatible(A, B))
	// Output: true <nil>
}

const simpleJSON = `{"firstname": "chris", "lastname": "james", "age": 30}`
const comparableJSON = `{"firstname": "christopher", "lastname": "james", "age": 15}`
const notSimilarJSON = `{"foo":"bar"}`

func TestItKnowsTheSameJSONIsCompatible(t *testing.T) {
	assertCompatible(t, simpleJSON, simpleJSON)
}

func TestItKnowsStructurallySameJSONIsCompatible(t *testing.T) {
	assertCompatible(t, simpleJSON, comparableJSON)
}

func TestItKnowsDifferentJSONIsIncompatible(t *testing.T) {
	assertIncompatible(t, simpleJSON, notSimilarJSON)
}

func TestItDoesntMindSuperflousFieldsInB(t *testing.T) {
	extraJSON := `{"firstname":"frank", "lastname": "sinatra", "extra field": "blue", "age":70}`
	assertCompatible(t, simpleJSON, extraJSON)
}

func TestItReturnsAnErrorForNonJson(t *testing.T) {
	if _, err := IsCompatible("nonsense", "not json"); err == nil {
		t.Error("Expected an error to be returned when both json is bad")
	}
	if _, err := IsCompatible(simpleJSON, "not json"); err == nil {
		t.Error("Expected an error to be returned when B is bad json")
	}
}
func TestFloatingPoints(t *testing.T) {
	floatingJSONa := `{"x": 3.14, "y": "not"}`
	floatingJSONb := `{"x": "three", "y": "not"}`
	assertIncompatible(t, floatingJSONa, floatingJSONb)
}

func TestStringsTypeCheck(t *testing.T) {
	stringyJSON := `{"x":"y"}`
	notStringyJSON := `{"x":1}`
	assertIncompatible(t, stringyJSON, notStringyJSON)
}

func TestBooleans(t *testing.T) {
	boolyJSONa := `{"x": true}`
	boolyJSONb := `{"x": false}`
	notBoolyJSON := `{"x": 1}`
	assertCompatible(t, boolyJSONa, boolyJSONb)
	assertIncompatible(t, boolyJSONa, notBoolyJSON)
}

func TestItKnowsHowToHandleSimpleArrays(t *testing.T) {
	JSONWithArray := `{"foo": ["baz", "bo"]}`
	comparableJSONWithArray := `{"foo": ["bar"]}`
	badlyTypedJSONArray := `{"foo": [1, 2]}`
	nonJSONArray := `{"foo":"bar"}`

	assertCompatible(t, JSONWithArray, comparableJSONWithArray)
	assertIncompatible(t, JSONWithArray, badlyTypedJSONArray)
	assertIncompatible(t, JSONWithArray, nonJSONArray)
}

func TestNestedStructures(t *testing.T) {
	a := `{"hello": [{"x": 1, "y": "a"},{"x": 2, "y": "b"}]}`
	b := `{"hello": [{"x": 10, "y": "b"}]}`
	c := `{"hello": [{"x": 10}]}`
	d := `{"hello": [{"z": 10}]}`
	e := `{"hello":[1,2,3]}`

	assertCompatible(t, a, b)
	assertIncompatible(t, a, c)
	assertIncompatible(t, a, d)
	assertIncompatible(t, a, e)
}

func TestEmptyArrayInB(t *testing.T) {
	a := `{"foo":["bar", "baz"]}`
	b := `{"foo":[]}`
	assertIncompatible(t, a, b)
}

func assertCompatible(t *testing.T, a, b string) {
	if compatible, err := IsCompatible(a, b); !compatible || err != nil {
		t.Errorf("%s should be compatible with %s (err = %v)", a, b, err)
	}
}

func assertIncompatible(t *testing.T, a, b string) {
	if compatible, err := IsCompatible(a, b); compatible || err != nil {
		t.Errorf("%s should not be compatible with %s (err = %v)", a, b, err)
	}
}

func BenchmarkItworks(b *testing.B) {
	a := `{"foo":[], "null": null}`
	// a := `{"total":0,"max_score":null,"hits":[],"categories":{"headlines":{"total":0,"labels":[]},"research":{"total":0,"labels":[]}}}`
	for i := 0; i < b.N; i++ {
		isCompat, err := IsCompatible(a, a)
		if err != nil {
			b.Fatal("IT ERRORD", err)
		}
		if !isCompat {
			b.Fatal("WHAT THE HELL MAN")
		}

	}
}
