package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for Codeforces problem described in problemE.txt (Price Maximization).
// We first count floor(a_i/k) for each item individually. The remainders are then
// paired greedily after sorting to gain additional revenue when their sum
// reaches at least k.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		var result int64
		rems := make([]int64, n)
		for i, v := range arr {
			result += v / k
			rems[i] = v % k
		}

		sort.Slice(rems, func(i, j int) bool { return rems[i] < rems[j] })
		i, j := 0, n-1
		for i < j {
			if rems[i]+rems[j] >= k {
				result++
				i++
				j--
			} else {
				i++
			}
		}

		fmt.Fprintln(writer, result)
	}
}
