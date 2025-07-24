package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var d int64
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}
	p := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}
	sort.Slice(p, func(i, j int) bool { return p[i] > p[j] })
	i, j := 0, n-1
	wins := 0
	for i <= j {
		needed := int(d/p[i]) + 1
		if i+needed-1 <= j {
			wins++
			i++
			j -= needed - 1
		} else {
			break
		}
	}
	fmt.Fprintln(out, wins)
}
