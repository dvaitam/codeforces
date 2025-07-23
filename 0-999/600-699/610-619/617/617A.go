package main

import "fmt"

func main() {
	var x int
	if _, err := fmt.Scan(&x); err != nil {
		return
	}
	// The elephant can move at most 5 units per step.
	// The minimal number of steps is ceil(x/5).
	steps := (x + 4) / 5
	fmt.Println(steps)
}
