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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	half := n / 2
	diff := 0
	q1, q2 := 0, 0
	for i := 0; i < half; i++ {
		if s[i] == '?' {
			q1++
		} else {
			diff += int(s[i] - '0')
		}
	}
	for i := half; i < n; i++ {
		if s[i] == '?' {
			q2++
		} else {
			diff -= int(s[i] - '0')
		}
	}
	if diff+(q1-q2)/2*9 == 0 {
		fmt.Fprintln(out, "Bicarp")
	} else {
		fmt.Fprintln(out, "Monocarp")
	}
}
