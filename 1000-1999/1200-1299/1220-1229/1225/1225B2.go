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
		var n, k, d int
		fmt.Fscan(reader, &n, &k, &d)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		freq := make(map[int]int)
		distinct := 0
		for i := 0; i < d; i++ {
			v := arr[i]
			if freq[v] == 0 {
				distinct++
			}
			freq[v]++
		}
		ans := distinct
		for i := d; i < n; i++ {
			rm := arr[i-d]
			freq[rm]--
			if freq[rm] == 0 {
				delete(freq, rm)
				distinct--
			}
			v := arr[i]
			if freq[v] == 0 {
				distinct++
			}
			freq[v]++
			if distinct < ans {
				ans = distinct
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
