package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// goodPair checks if two activity sets form a valid match.
func goodPair(a, b []int) bool {
	i, j := 0, 0
	shared := false
	diffA, diffB := false, false
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			shared = true
			i++
			j++
		} else if a[i] < b[j] {
			diffA = true
			i++
		} else {
			diffB = true
			j++
		}
		if shared && diffA && diffB {
			return true
		}
	}
	if i < len(a) {
		diffA = true
	}
	if j < len(b) {
		diffB = true
	}
	return shared && diffA && diffB
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	acts := make([][]int, n)
	actUsers := make(map[int][]int)
	for i := 0; i < n; i++ {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return
		}
		arr := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(reader, &arr[j])
		}
		sort.Ints(arr)
		acts[i] = arr
		for _, x := range arr {
			actUsers[x] = append(actUsers[x], i)
		}
	}

	threshold := 500
	visited := make(map[uint64]struct{})

	for _, list := range actUsers {
		if len(list) <= 1 {
			continue
		}
		if len(list) <= threshold {
			for i := 0; i < len(list); i++ {
				u := list[i]
				if len(acts[u]) < 2 {
					continue
				}
				for j := i + 1; j < len(list); j++ {
					v := list[j]
					if len(acts[v]) < 2 {
						continue
					}
					var key uint64
					if u < v {
						key = uint64(u)<<32 | uint64(v)
					} else {
						key = uint64(v)<<32 | uint64(u)
					}
					if _, ok := visited[key]; ok {
						continue
					}
					visited[key] = struct{}{}
					if goodPair(acts[u], acts[v]) {
						fmt.Fprintln(writer, "YES")
						if u < v {
							fmt.Fprintf(writer, "%d %d\n", u+1, v+1)
						} else {
							fmt.Fprintf(writer, "%d %d\n", v+1, u+1)
						}
						return
					}
				}
			}
		} else {
			// For large lists, check pairs involving the first 'threshold' users only.
			limit := threshold
			if len(list) < limit {
				limit = len(list)
			}
			for idx := 0; idx < limit; idx++ {
				u := list[idx]
				if len(acts[u]) < 2 {
					continue
				}
				for j := idx + 1; j < len(list); j++ {
					v := list[j]
					if len(acts[v]) < 2 {
						continue
					}
					var key uint64
					if u < v {
						key = uint64(u)<<32 | uint64(v)
					} else {
						key = uint64(v)<<32 | uint64(u)
					}
					if _, ok := visited[key]; ok {
						continue
					}
					visited[key] = struct{}{}
					if goodPair(acts[u], acts[v]) {
						fmt.Fprintln(writer, "YES")
						if u < v {
							fmt.Fprintf(writer, "%d %d\n", u+1, v+1)
						} else {
							fmt.Fprintf(writer, "%d %d\n", v+1, u+1)
						}
						return
					}
				}
			}
		}
	}

	fmt.Fprintln(writer, "NO")
}
