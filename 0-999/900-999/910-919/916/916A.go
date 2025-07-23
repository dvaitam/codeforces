package main

import (
	"bufio"
	"fmt"
	"os"
)

func contains7(n int) bool {
	return n%10 == 7 || n/10 == 7
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var x, hh, mm int
	if _, err := fmt.Fscan(in, &x, &hh, &mm); err != nil {
		return
	}

	cnt := 0
	for {
		if contains7(hh) || contains7(mm) {
			fmt.Fprintln(out, cnt)
			return
		}
		mm -= x
		for mm < 0 {
			mm += 60
			hh--
		}
		if hh < 0 {
			hh += 24
		}
		cnt++
	}
}
