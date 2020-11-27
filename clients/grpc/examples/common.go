package main

import (
	"log"
	"reflect"
)

func AssertEqual(expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) == false {
		log.Fatalf("Expected: %v Actual: %v\n", expected, actual)
	}
}

func AssertNil(object interface{}) {
	if !reflect.ValueOf(object).IsNil() {
		log.Fatalf("%v must be nil.\n", object)
	}
}

func AssertTrue(actual interface{}) {
	if actual == false {
		log.Fatalf("must be true.\n")
	}
}
