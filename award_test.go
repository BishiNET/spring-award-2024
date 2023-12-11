package newaward

import "testing"

func TestNewAward(t *testing.T) {
	a := NewAward(
		WithOneAction(func() {
			t.Log("big award")
		}),
		WithLeftChildAction(func() {
			t.Log("left big award")
		}, func() {
			t.Log("left 20%")
		}, func() {
			t.Log("left 79%")
		}),
		WithRightChildAction(func() {
			t.Log("left big award")
		}, func() {
			t.Log("right 20%")
		}, func() {
			t.Log("right 79%")
		}),
	)

	for i := 0; i < 10; i++ {
		a.Pick()
	}
}
