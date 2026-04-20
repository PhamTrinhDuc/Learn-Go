package main

import "fmt"

type Student struct {
	Id     int
	Name   string
	Grades []float64
}

func (s Student) Averge() float64 {
	sum := 0.0
	for _, grade := range s.Grades {
		sum += grade
	}
	return sum / float64(len(s.Grades))
}

func (s Student) IsPassed() bool {
	if s.Averge() >= 5 {
		for _, x := range s.Grades {
			if x < 3 {
				return false
			}
		}
		return true
	}
	return false
}

type ClassRoom struct {
	Name     string
	Students []Student
}

func (c ClassRoom) GetTopStudent() Student {
	student := Student{}
	avg_highest := 0.0
	for _, st := range c.Students {
		if st.Averge() > avg_highest {
			avg_highest = st.Averge()
			student = st
		}
	}
	return student
}

func main() {
	student := Student{
		Id:     1,
		Name:   "John",
		Grades: []float64{10, 9, 8},
	}
	student2 := Student{
		Id:     2,
		Name:   "Jane",
		Grades: []float64{9, 8, 7},
	}
	student3 := Student{
		Id:     3,
		Name:   "Bob",
		Grades: []float64{8, 7, 6},
	}
	classroom := ClassRoom{
		Name:     "Classroom 1",
		Students: []Student{student, student2, student3},
	}
	fmt.Println(classroom.GetTopStudent())
}
