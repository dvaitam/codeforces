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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g := make([]string, n)
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(reader, &line)
		g[i] = line
	}
	type pair struct{ val, idx int }
	v1 := make([]pair, n)
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < m; j++ {
			if g[i][j] == '>' {
				cnt++
			} else if g[i][j] == '<' {
				cnt--
			}
		}
		v1[i] = pair{cnt, i}
	}
	sort.Slice(v1, func(i, j int) bool {
		return v1[i].val < v1[j].val
	})
	a := make([]int, n)
	a[v1[0].idx] = 1
	for i := 1; i < n; i++ {
		u := v1[i].idx
		v := v1[i-1].idx
		if v1[i].val == v1[i-1].val {
			a[u] = a[v]
		} else {
			a[u] = a[v] + 2
		}
	}
	b := make([]int, m)
	for j := 0; j < m; j++ {
		assigned := false
		for i := 0; i < n; i++ {
			if g[i][j] == '=' {
				b[j] = a[i]
				assigned = true
				break
			}
		}
		if !assigned {
			for i := 0; i < n; i++ {
				if g[i][j] == '<' && a[i]+1 > b[j] {
					b[j] = a[i] + 1
				}
			}
		}
	}
	ok := true
	for i := 0; i < n && ok; i++ {
		for j := 0; j < m; j++ {
			var tmp byte
			if a[i] > b[j] {
				tmp = '>'
			} else if a[i] == b[j] {
				tmp = '='
			} else {
				tmp = '<'
			}
			if tmp != g[i][j] {
				ok = false
				break
			}
		}
	}
	if !ok {
		fmt.Fprintln(writer, "No")
		return
	}
	fmt.Fprintln(writer, "Yes")
	ve := make([]int, 0, n+m)
	for i := 0; i < n; i++ {
		ve = append(ve, a[i])
	}
	for j := 0; j < m; j++ {
		ve = append(ve, b[j])
	}
	sort.Ints(ve)
	uniq := ve[:0]
	for i, v := range ve {
		if i == 0 || v != ve[i-1] {
			uniq = append(uniq, v)
		}
	}
	ve = uniq
	for i := 0; i < n; i++ {
		idx := sort.SearchInts(ve, a[i])
		a[i] = idx + 1
	}
	for j := 0; j < m; j++ {
		idx := sort.SearchInts(ve, b[j])
		b[j] = idx + 1
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(writer, "%d ", a[i])
	}
	fmt.Fprintln(writer)
	for j := 0; j < m; j++ {
		fmt.Fprintf(writer, "%d ", b[j])
	}
	fmt.Fprintln(writer)
}
