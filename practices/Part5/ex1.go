package main

import "fmt"

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	r := Rectangle{Width: 30, Height: 40}
	area := r.Area()
	peri := r.Perimeter()
	fmt.Println("Area: ", area)
	fmt.Println("Perimeter: ", peri)
	r.Scale(2)
	fmt.Println("Width: ", r.Width)
	fmt.Println("Height: ", r.Height)
}
