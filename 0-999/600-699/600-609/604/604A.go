package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It computes Kevin's final score based on submission times,
// wrong submissions and hacks according to Codeforces rules.
func main() {
	in := bufio.NewReader(os.Stdin)
	m := make([]int, 5)
	w := make([]int, 5)
	for i := 0; i < 5; i++ {
		if _, err := fmt.Fscan(in, &m[i]); err != nil {
			return
		}
	}
	for i := 0; i < 5; i++ {
		if _, err := fmt.Fscan(in, &w[i]); err != nil {
			return
		}
	}
	var hs, hu int
	if _, err := fmt.Fscan(in, &hs, &hu); err != nil {
		return
	}

	x := []int{500, 1000, 1500, 2000, 2500}
	total := 0
	for i := 0; i < 5; i++ {
		score1 := 3 * x[i] / 10
		score2 := (250-m[i])*x[i]/250 - 50*w[i]
		if score1 > score2 {
			total += score1
		} else {
			total += score2
		}
	}
	total += hs * 100
	total -= hu * 50

	fmt.Println(total)
}
