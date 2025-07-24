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

	a := make([]int, 6)
	sum := 0
	for i := 0; i < 6; i++ {
		if _, err := fmt.Fscan(in, &a[i]); err != nil {
			return
		}
		sum += a[i]
	}

	if sum%2 != 0 {
		fmt.Fprintln(out, "NO")
		return
	}
	target := sum / 2

	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			for k := j + 1; k < 6; k++ {
				if a[i]+a[j]+a[k] == target {
					fmt.Fprintln(out, "YES")
					return
				}
			}
		}
	}
	fmt.Fprintln(out, "NO")
}
