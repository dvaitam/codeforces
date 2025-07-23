package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	zeroIdx := -1
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == 0 {
			zeroIdx = i
		}
	}
	b := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &b[i])
	}

	if k > 1 {
		fmt.Fprintln(writer, "Yes")
		return
	}

	// k == 1, only one position to fill
	if zeroIdx != -1 {
		a[zeroIdx] = b[0]
	}

	increasing := true
	for i := 1; i < n; i++ {
		if a[i] <= a[i-1] {
			increasing = false
			break
		}
	}

	if increasing {
		fmt.Fprintln(writer, "No")
	} else {
		fmt.Fprintln(writer, "Yes")
	}
}
