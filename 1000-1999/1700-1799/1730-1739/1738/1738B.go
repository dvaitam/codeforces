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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		s := make([]int64, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &s[i])
		}
		if k == 1 {
			fmt.Fprintln(writer, "YES")
			continue
		}

		diff := make([]int64, k-1)
		for i := 1; i < k; i++ {
			diff[i-1] = s[i] - s[i-1]
		}
		ok := true
		for i := 1; i < len(diff); i++ {
			if diff[i] < diff[i-1] {
				ok = false
				break
			}
		}
		if ok {
			firstCount := int64(n - k + 1)
			if diff[0]*firstCount < s[0] {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
