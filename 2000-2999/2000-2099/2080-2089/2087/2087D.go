package main

import (
	"bufio"
	"fmt"
	"os"
)

func isUpper(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

func compareTarget(target, s string, upperFirst bool) int {
	th := isUpper(target[0])
	ts := isUpper(s[0])
	if th == ts {
		if target == s {
			return 0
		}
		if target < s {
			return -1
		}
		return 1
	}
	if upperFirst {
		if th {
			return -1
		}
		return 1
	}
	if th {
		return 1
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var target string
	if _, err := fmt.Fscan(in, &n, &target); err != nil {
		return
	}

	query := func(idx int) string {
		fmt.Fprintf(out, "? %d\n", idx)
		out.Flush()
		var resp string
		if _, err := fmt.Fscan(in, &resp); err != nil {
			return ""
		}
		return resp
	}

	first := query(1)
	if first == "" || first == "-1" {
		return
	}
	upperFirst := isUpper(first[0])

	if compareTarget(target, first, upperFirst) == 0 {
		fmt.Fprintln(out, "! 1")
		out.Flush()
		return
	}

	left, right := 2, n
	for left <= right {
		mid := (left + right) / 2
		cur := query(mid)
		if cur == "" || cur == "-1" {
			return
		}
		cmp := compareTarget(target, cur, upperFirst)
		if cmp == 0 {
			fmt.Fprintf(out, "! %d\n", mid)
			out.Flush()
			return
		}
		if cmp > 0 {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
}
