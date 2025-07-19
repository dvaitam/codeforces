package main

import "fmt"

func main() {
	// Total number of segments
	fmt.Println(2001)
	// First batch of segments
	for i := 1; i <= 1000; i++ {
		fmt.Printf("%d 1 %d 2\n", i, i)
	}
	// Middle segment
	fmt.Printf("1000 1 1000 2\n")
	// Second batch of segments
	for i := 1; i <= 1000; i++ {
		fmt.Printf("%d 1 %d 2\n", i, i)
	}
}
