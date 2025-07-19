package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	id  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)

	s := make([]pair, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		s[i] = pair{val: x, id: i + 1}
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].val < s[j].val
	})

	v := make([]pair, m)
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(reader, &x)
		v[i] = pair{val: x, id: i + 1}
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].val < v[j].val
	})

	usedN := make([]bool, n)
	usedM := make([]bool, m)
	ans := make([]int, n+1)
	taim := make([]int, m+1)

	cnt, sumT := 0, 0
	// Try up to 32 halvings
	for t := 0; t < 32; t++ {
		j := 0
		for i := 0; i < m; i++ {
			if usedM[i] {
				// still halve for consistency
				v[i].val = (v[i].val >> 1) + (v[i].val & 1)
				continue
			}
			val := v[i].val
			// find first unused s[j] with s[j].val >= val
			for j < n && (usedN[j] || s[j].val < val) {
				j++
			}
			if j < n && s[j].val == val {
				// match
				cnt++
				sumT += t
				ans[s[j].id] = v[i].id
				taim[v[i].id] = t
				usedN[j] = true
				usedM[i] = true
			}
			// halve the value (ceil)
			v[i].val = (v[i].val >> 1) + (v[i].val & 1)
		}
	}

	// Output results
	fmt.Fprintln(writer, cnt, sumT)
	for id := 1; id <= m; id++ {
		if id > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, taim[id])
	}
	fmt.Fprintln(writer)
	for id := 1; id <= n; id++ {
		if id > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[id])
	}
	fmt.Fprintln(writer)
}
