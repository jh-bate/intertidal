package intertidal

import (
	"testing"
)

func Test_meetTarget(t *testing.T) {

	p1 := &Pledge{TargetValue: "> 3.5"}

	if p1.meetTarget(9.0) != true {
		t.Fatal("opps should be true")
	}

	p2 := &Pledge{TargetValue: "3"}

	if p2.meetTarget(3.0) != true {
		t.Fatal("opps should be true")
	}

	p3 := &Pledge{TargetValue: "<= 9"}

	if p3.meetTarget(9.1) != false {
		t.Fatal("opps should be false")
	}

}
