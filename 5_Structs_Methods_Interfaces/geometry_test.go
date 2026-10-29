package main

import "testing"

func TestPerimeter(t *testing.T) {
	got := Perimeter(Rectangle{4.0, 5.0})
	want := 18.0

	if got != want {
		t.Errorf("Got: %.2f, Want: %.2f", got, want)
	}
}

/*
	=>	But what if a bug slipped in to the system in one of 20 cases in the table? How would a developer know which case failed? This is not a great experience for the
		developer, they will have to manually look through the cases to find out which case actually failed. We can change our error message into %#v got %g want %g.
	=>	The %#v format string will print out our struct with the values in its field, so the developer can see at a glance the properties that are being tested.
	=>	One final tip with table driven tests is to use t.Run and to name the test cases. By wrapping each case in a t.Run you will have clearer test output on failures
		as it will print the name of the case.
*/

func TestArea(t *testing.T) {

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Length: 12, Breadth: 6}, hasArea: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
			}
		})
	}
}
