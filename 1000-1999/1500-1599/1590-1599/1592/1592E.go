package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct {
	val   int
	start int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	pref := make([]int, n+1)
	prev := []Pair{}
	ans := 0

	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ arr[i]
		cur := []Pair{{arr[i], i}}
		for _, p := range prev {
			val := p.val & arr[i]
			if cur[len(cur)-1].val == val {
				if p.start < cur[len(cur)-1].start {
					cur[len(cur)-1].start = p.start
				}
			} else {
				cur = append(cur, Pair{val, p.start})
			}
		}
		for _, p := range cur {
			xorVal := pref[i] ^ pref[p.start-1]
			if p.val > xorVal {
				length := i - p.start + 1
				if length > ans {
					ans = length
				}
			}
		}
		prev = cur
	}

	fmt.Fprintln(writer, ans)
}
