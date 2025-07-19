package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var tt int
	fmt.Fscan(reader, &tt)
	for tt > 0 {
		tt--
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	minv := math.MaxInt64
	s := make(map[int]struct{})
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] < minv {
			minv = a[i]
		}
		s[a[i]] = struct{}{}
	}
	if minv < 0 {
		fmt.Fprintln(writer, "NO")
		return
	}
	fmt.Fprintln(writer, "YES")
	sort.Ints(a)
	// closure under absolute difference
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			tmp := a[j] - a[i]
			if tmp < 0 {
				tmp = -tmp
			}
			if _, exists := s[tmp]; !exists {
				s[tmp] = struct{}{}
				a = append(a, tmp)
				sort.Ints(a)
				i = -1
				break
			}
		}
	}
	// prepare and print result
	res := make([]int, 0, len(s))
	for v := range s {
		res = append(res, v)
	}
	sort.Ints(res)
	fmt.Fprintln(writer, len(res))
	for i, v := range res {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
