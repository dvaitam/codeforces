package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: naive recursive solver for demonstration; does not scale to limits.

var memo map[[2]int64]bool

func win(s, e int64) bool {
	if s > e {
		return false
	}
	key := [2]int64{s, e}
	if v, ok := memo[key]; ok {
		return v
	}
	var res bool
	if 2*s > e {
		res = (e-s)%2 == 1
	} else {
		res = !(win(s+1, e) && win(2*s, e))
	}
	memo[key] = res
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	memo = make(map[[2]int64]bool)

	var s, e int64
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &s, &e)
		if win(s, e) {
			fmt.Fprint(writer, "1 0\n")
		} else {
			fmt.Fprint(writer, "0 1\n")
		}
	}
}
