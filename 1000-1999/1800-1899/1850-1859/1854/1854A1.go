package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ x, y int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for tc := 0; tc < T; tc++ {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	mx, mn := 0, 0
	for i := 0; i < n; i++ {
		if a[i] > a[mx] {
			mx = i
		}
		if a[i] < a[mn] {
			mn = i
		}
	}

	sorted := true
	for i := 1; i < n; i++ {
		if a[i] < a[i-1] {
			sorted = false
			break
		}
	}
	if sorted {
		fmt.Fprintln(writer, 0)
		return
	}

	br := make([]pair, 0)
	cr := make([]pair, 0)

	if a[mx] > 0 {
		boro := a[mx]
		for boro+a[mn] < 0 {
			br = append(br, pair{mx + 1, mx + 1})
			boro *= 2
		}
		for i := 1; i < n; i++ {
			if a[i] < 0 {
				br = append(br, pair{i + 1, mx + 1})
			}
		}
		for i := 1; i < n; i++ {
			br = append(br, pair{i + 1, i})
		}
	}

	if a[mn] < 0 {
		choto := a[mn]
		for choto+a[mx] > 0 {
			cr = append(cr, pair{mn + 1, mn + 1})
			choto *= 2
		}
		for i := n - 2; i >= 0; i-- {
			if a[i] > 0 {
				cr = append(cr, pair{i + 1, mn + 1})
			}
		}
		for i := n - 1; i > 0; i-- {
			cr = append(cr, pair{i, i + 1})
		}
	}

	var ans []pair
	if len(br) < len(cr) {
		ans = br
	} else {
		ans = cr
	}

	fmt.Fprintln(writer, len(ans))
	for _, p := range ans {
		fmt.Fprintln(writer, p.x, p.y)
	}
}
