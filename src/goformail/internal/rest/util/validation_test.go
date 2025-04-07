package util

import (
	"testing"
)

func TestValidateAllSet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{String: "a", Number: 1, Array: []string{}})

	if !actual {
		t.Fatalf("Validator returned false when all fields were set")
	}
}

func TestValidateNoStringSet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{Number: 1, Array: []string{}})

	if actual {
		t.Fatalf("Validator returned true when no string was set")
	}
}

func TestValidateNoIntSet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{String: "a", Array: []string{}})

	if actual {
		t.Fatalf("Validator returned true when no int was set")
	}
}

func TestValidateNoArraySet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{String: "a", Number: 1})

	if actual {
		t.Fatalf("Validator returned true when no array was set")
	}
}

type TestStruct struct {
	String string
	Number int
	Array  []string
}
