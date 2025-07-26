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
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		cntA := make(map[int]int)
		cntB := make(map[int]int)

		for i := 0; i < n; {
			j := i
			for j < n && a[j] == a[i] {
				j++
			}
			length := j - i
			if length > cntA[a[i]] {
				cntA[a[i]] = length
			}
			i = j
		}

		for i := 0; i < n; {
			j := i
			for j < n && b[j] == b[i] {
				j++
			}
			length := j - i
			if length > cntB[b[i]] {
				cntB[b[i]] = length
			}
			i = j
		}

		ans := 0
		for v, la := range cntA {
			if la+cntB[v] > ans {
				ans = la + cntB[v]
			}
		}
		for v, lb := range cntB {
			if lb > ans && cntA[v] == 0 {
				ans = lb
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
