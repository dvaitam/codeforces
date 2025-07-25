package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		found := false
		for i := 1; i < len(s); i++ {
			aStr := s[:i]
			bStr := s[i:]
			if aStr[0] == '0' || bStr[0] == '0' {
				continue
			}
			aVal, _ := strconv.Atoi(aStr)
			bVal, _ := strconv.Atoi(bStr)
			if bVal > aVal {
				fmt.Fprintf(out, "%d %d\n", aVal, bVal)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintln(out, -1)
		}
	}
}
