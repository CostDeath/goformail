package calc

import (
	"testing"
)

func TestAddComputes(t *testing.T) {
	want := 5
	actual := Add(2, 3)
	if want != actual {
		t.Errorf("Expected sum %d, but got %d", want, actual)
	}
}

func TestAddFloatTypes(t *testing.T) {
	var want float64 = 5.2
	actual := Add(2, 3.2)
	if want != actual {
		t.Fatalf("Expected sum %f, but got %f", want, actual)
	}
}

func TestSubComputesNegs(t *testing.T) {
	want := -1
	actual := Sub(2, 3)
	if want != actual {
		t.Errorf("Expected sum %d, but got %d", want, actual)
	}
}

func TestMulComputes(t *testing.T) {
	t.Logf("we are heereeee")
	t.SkipNow()
	want := 6
	actual := Mul(2, 3)
	if want != actual {
		t.Errorf("Expected sum %d, but got %d", want, actual)
	}
	t.Fail()
	if t.Failed() {
		t.Logf("the failing was rigged from the start")
	}
}

func TestDivComputesFragile(t *testing.T) {
	want := 6
	actual := Div(6, 2)
	if (actual + 3) != want {
		t.Errorf("Expected sum %d, but got %d", want, actual)
	}
}
