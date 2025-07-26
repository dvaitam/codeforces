package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

// digits of Pi used to determine the length of each test case
var piDigits = []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9, 3, 2, 3, 8, 4, 6, 2, 6, 4, 3, 3, 8, 3, 2, 7, 9, 5}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		n := piDigits[i]
		prod := big.NewInt(1)
		for j := 0; j < n; j++ {
			var x int64
			fmt.Fscan(in, &x)
			prod.Mul(prod, big.NewInt(x))
		}
		fmt.Println(prod)
	}
}
