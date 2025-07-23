package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var a, b int64
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}

	var count int64
	for l := 2; l <= 60; l++ {
		allOnes := (int64(1) << uint(l)) - 1
		for i := 0; i < l-1; i++ {
			num := allOnes - (int64(1) << uint(i))
			if num >= a && num <= b {
				count++
			}
		}
	}

	fmt.Fprintln(out, count)
}
