package main

import "fmt"

type mt map[int]int

func main() {
	mm := make(mt)
	loadmap(mm)
	checkmap(mm)
	fmt.Println("Hey it worked")
}

func loadmap(m mt) {
	for i := range 1000 {
		m[i] = i + 100
	}
}

func checkmap(m mt) {
	for i := range 1000 {
		if m[i] != i+100 {
			panic(i)
		}
	}
}
