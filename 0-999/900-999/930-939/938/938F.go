package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	n := len(s)
	if n == 0 {
		fmt.Fprintln(writer, "")
		return
	}
	k := bits.Len(uint(n)) - 1
	lengths := make([]int, k)
	for i := 0; i < k; i++ {
		lengths[i] = 1 << i
	}
	finalLen := n - ((1 << k) - 1)

	state := func(pos, mask int) int { return mask*(n+1) + pos }

	cur := map[int]struct{}{state(0, 0): {}}
	ans := make([]byte, 0, finalLen)

	for iter := 0; iter < finalLen; iter++ {
		// expand all possible deletions using BFS
		queue := make([]int, 0, len(cur))
		visited := make(map[int]struct{}, len(cur))
		for st := range cur {
			visited[st] = struct{}{}
			queue = append(queue, st)
		}
		for len(queue) > 0 {
			st := queue[0]
			queue = queue[1:]
			pos := st % (n + 1)
			mask := st / (n + 1)
			for i := 0; i < k; i++ {
				if (mask>>i)&1 == 1 {
					continue
				}
				l := lengths[i]
				if pos+l <= n {
					ns := state(pos+l, mask|1<<i)
					if _, ok := visited[ns]; !ok {
						visited[ns] = struct{}{}
						queue = append(queue, ns)
					}
				}
			}
		}

		// choose smallest possible next character
		minc := byte('{')
		for st := range visited {
			pos := st % (n + 1)
			if pos < n {
				c := s[pos]
				if c < minc {
					minc = c
				}
			}
		}
		ans = append(ans, minc)

		// prepare states for next position
		next := make(map[int]struct{})
		for st := range visited {
			pos := st % (n + 1)
			mask := st / (n + 1)
			if pos < n && s[pos] == minc {
				next[state(pos+1, mask)] = struct{}{}
			}
		}
		cur = next
	}

	fmt.Fprintln(writer, string(ans))
}
