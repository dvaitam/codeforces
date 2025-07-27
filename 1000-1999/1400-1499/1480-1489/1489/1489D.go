package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// This program solves the problem described in problemD.txt.
// We are given n strings and must reorder them so that each string
// contains all previous strings as substrings. The easiest approach
// is to sort the strings by length in nondecreasing order. If such
// ordering is valid, every shorter string will appear in every longer
// string consecutively; otherwise it's impossible.
// After sorting, we simply verify that for every i, strings[i] is a
// substring of strings[i+1]. If any check fails we output "NO".
// Otherwise we output "YES" followed by the ordered list.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	strs := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &strs[i])
	}

	sort.Slice(strs, func(i, j int) bool {
		if len(strs[i]) == len(strs[j]) {
			return strs[i] < strs[j]
		}
		return len(strs[i]) < len(strs[j])
	})

	for i := 0; i < n-1; i++ {
		if !strings.Contains(strs[i+1], strs[i]) {
			fmt.Fprintln(out, "NO")
			return
		}
	}

	fmt.Fprintln(out, "YES")
	for _, s := range strs {
		fmt.Fprintln(out, s)
	}
}
