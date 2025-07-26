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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, l, r int64
		fmt.Fscan(reader, &n, &l, &r)
		res := solve(n, l, r)
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}

func solve(n, l, r int64) []int64 {
	result := make([]int64, 0, r-l+1)
	var pos int64 = 1
	for i := int64(1); i < n && pos <= r; i++ {
		blockLen := 2 * (n - i)
		if l > pos+blockLen-1 {
			pos += blockLen
			continue
		}
		for j := i + 1; j <= n && pos <= r; j++ {
			if pos >= l {
				result = append(result, i)
			}
			pos++
			if pos > r {
				break
			}
			if pos >= l {
				result = append(result, j)
			}
			pos++
		}
	}
	if pos <= r {
		// last element is always 1
		if pos >= l {
			result = append(result, 1)
		}
	}
	return result
}
