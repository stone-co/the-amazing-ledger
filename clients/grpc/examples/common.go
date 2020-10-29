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
