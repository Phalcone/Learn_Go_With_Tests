package main

import "math"

/*
	We're creating a new type just like we did with Rectangle and Circle but this time it is an interface rather than a struct.
	Rectangle has a method called Area that returns a float64 so it satisfies the Shape interface
	Circle has a method called Area that returns a float64 so it satisfies the Shape interface
	String does not have such a method, so it doesn't satisfy the interface
*/

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Length, Breadth float64
}

func (rectangle Rectangle) Area() float64 {
	return rectangle.Breadth * rectangle.Length
}

type Circle struct {
	Radius float64
}

func (circle Circle) Area() float64 {
	return math.Pi * (math.Pow(circle.Radius, 2.0))
}

type Triangle struct {
	Base, Height float64
}

func (triangle Triangle) Area() float64 {
	return 0.5 * triangle.Base * triangle.Height
}

func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.Length + rectangle.Breadth)
}
