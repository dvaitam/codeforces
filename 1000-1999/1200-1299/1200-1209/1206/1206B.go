package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	var cost int64
	negatives := 0
	zeros := 0
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		if x > 0 {
			cost += x - 1
		} else if x < 0 {
			cost += -1 - x
			negatives++
		} else {
			zeros++
		}
	}

	if negatives%2 == 0 {
		cost += int64(zeros)
	} else {
		if zeros > 0 {
			cost += int64(zeros)
		} else {
			cost += 2
		}
	}

	fmt.Println(cost)
}
