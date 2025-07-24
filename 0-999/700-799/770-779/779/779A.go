package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	freqA := make([]int, 6)
	freqB := make([]int, 6)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		freqA[x]++
	}
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		freqB[x]++
	}
	sumDiff := 0
	for i := 1; i <= 5; i++ {
		if (freqA[i]+freqB[i])%2 != 0 {
			fmt.Println(-1)
			return
		}
		d := freqA[i] - freqB[i]
		if d < 0 {
			d = -d
		}
		sumDiff += d
	}
	fmt.Println(sumDiff / 4)
}
