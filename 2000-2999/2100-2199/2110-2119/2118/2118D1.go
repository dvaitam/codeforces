package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const limit int64 = 1_000_000_000_000_000 // 1e15

type state struct {
	idx int
	dir int8
	t   int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		pos := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &pos[i])
		}
		delay := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &delay[i])
		}
		posToIdx := make(map[int64]int, n)
		for i, p := range pos {
			posToIdx[p] = i
		}
		var q int
		fmt.Fscan(reader, &q)
		starts := make([]int64, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(reader, &starts[i])
		}
		for _, start := range starts {
			dir := int8(1)
			var tmod int64
			cur := start
			visited := make(map[state]struct{})
			exit := false
			for {
				if cur < 1 || cur > limit {
					exit = true
					break
				}
				if idx, ok := posToIdx[cur]; ok {
					st := state{idx: idx, dir: dir, t: int(tmod)}
					if _, seen := visited[st]; seen {
						break
					}
					visited[st] = struct{}{}
					if tmod == int64(delay[idx]) {
						dir = -dir
					}
				}
				var steps int64
				if dir == 1 {
					j := sort.Search(len(pos), func(i int) bool { return pos[i] > cur })
					if j == len(pos) {
						steps = limit + 1 - cur
						cur += steps
						if cur > limit {
							exit = true
						}
						break
					}
					nextPos := pos[j]
					steps = nextPos - cur
					cur = nextPos
				} else {
					j := sort.Search(len(pos), func(i int) bool { return pos[i] >= cur })
					if j == 0 {
						steps = cur
						cur -= steps
						if cur < 1 {
							exit = true
						}
						break
					}
					prevPos := pos[j-1]
					steps = cur - prevPos
					cur = prevPos
				}
				tmod = (tmod + steps%int64(k)) % int64(k)
			}
			if exit {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
