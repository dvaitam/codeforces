package main

import "fmt"

func main() {
	// TODO: Implement solution for problem F - Intergalactic Sliding Puzzle
	var t int
	fmt.Scan(&t)
	for ; t > 0; t-- {
		var k int
		fmt.Scan(&k)
		// Read the two rows but ignore values for now
		for i := 0; i < 2; i++ {
			for j := 0; j < 2*k+1; j++ {
				var tmp string
				fmt.Scan(&tmp)
			}
		}
		fmt.Println("SURGERY FAILED")
	}
}
