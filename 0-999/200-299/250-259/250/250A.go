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

	days := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &days[i])
	}

	var folders []int
	currentLen := 0
	negatives := 0

	for _, val := range days {
		currentLen++
		if val < 0 {
			negatives++
		}

		if negatives == 3 {
			// Close the previous folder before this third negative.
			folders = append(folders, currentLen-1)
			currentLen = 1                 // current folder starts with current day
			negatives = boolToInt(val < 0) // current folder negative count
		}
	}

	folders = append(folders, currentLen)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fprintln(out, len(folders))
	for i, v := range folders {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}

func boolToInt(cond bool) int {
	if cond {
		return 1
	}
	return 0
}
