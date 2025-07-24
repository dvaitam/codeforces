package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}
	seq := []int{4, 22, 27, 58, 85, 94, 121, 166, 202, 265, 274, 319, 346, 355, 378, 382, 391, 438, 454, 483, 517, 526, 535, 562, 576, 588, 627, 634, 636, 645}
	if a >= 1 && a <= len(seq) {
		fmt.Println(seq[a-1])
	}
}
