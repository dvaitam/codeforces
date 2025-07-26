package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		type pair struct{ req, gain int }
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{a[i], b[i]}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].req < arr[j].req })
		ram := k
		for _, p := range arr {
			if ram < p.req {
				break
			}
			ram += p.gain
		}
		fmt.Fprintln(writer, ram)
	}
}
