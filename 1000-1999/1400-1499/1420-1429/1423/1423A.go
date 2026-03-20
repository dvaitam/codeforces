package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cost struct {
	city int
	cost int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	if !scanner.Scan() {
		return
	}

	n := 0
	for _, v := range scanner.Bytes() {
		n = n*10 + int(v-'0')
	}

	if n%2 != 0 {
		fmt.Println("-1")
		return
	}

	pref := make([][]int, n+1)
	rank := make([][]int, n+1)

	for i := 1; i <= n; i++ {
		pref[i] = make([]int, n-1)
		rank[i] = make([]int, n+1)

		costs := make([]Cost, n-1)
		idx := 0
		for j := 1; j <= n; j++ {
			if i == j {
				continue
			}
			scanner.Scan()
			c := 0
			for _, v := range scanner.Bytes() {
				c = c*10 + int(v-'0')
			}
			costs[idx] = Cost{city: j, cost: c}
			idx++
		}

		sort.Slice(costs, func(a, b int) bool {
			return costs[a].cost < costs[b].cost
		})

		for j := 0; j < n-1; j++ {
			pref[i][j] = costs[j].city
			rank[i][costs[j].city] = j
		}
	}

	head := make([]int, n+1)
	tail := make([]int, n+1)
	next_val := make([][]int, n+1)
	prev_val := make([][]int, n+1)
	match := make([]int, n+1)

	for i := 1; i <= n; i++ {
		head[i] = 0
		tail[i] = n - 2
		next_val[i] = make([]int, n-1)
		prev_val[i] = make([]int, n-1)
		for j := 0; j < n-1; j++ {
			next_val[i][j] = j + 1
			prev_val[i][j] = j - 1
		}
		next_val[i][n-2] = -1
	}

	remove := func(i, idx int) {
		p := prev_val[i][idx]
		nx := next_val[i][idx]
		if p != -1 {
			next_val[i][p] = nx
		} else {
			head[i] = nx
		}
		if nx != -1 {
			prev_val[i][nx] = p
		} else {
			tail[i] = p
		}
	}

	Q := make([]int, 0, n*n)
	for i := 1; i <= n; i++ {
		Q = append(Q, i)
	}

	processQ := func() bool {
		for len(Q) > 0 {
			x := Q[0]
			Q = Q[1:]

			if head[x] == -1 {
				return false
			}

			y := pref[x][head[x]]
			idx_y := rank[y][x]

			curr := next_val[y][idx_y]
			for curr != -1 {
				w := pref[y][curr]
				remove(w, rank[w][y])
				if head[w] == -1 {
					return false
				}
				if match[y] == w {
					match[y] = 0
					Q = append(Q, w)
				}
				curr = next_val[y][curr]
			}

			tail[y] = idx_y
			next_val[y][idx_y] = -1
			match[y] = x
		}
		return true
	}

	if !processQ() {
		fmt.Println("-1")
		return
	}

	in_path := make([]int, n+1)
	for i := 0; i <= n; i++ {
		in_path[i] = -1
	}

	for {
		found_cycle := false
		for i := 1; i <= n; i++ {
			if head[i] != -1 && head[i] != tail[i] {
				found_cycle = true

				seq := []int{}
				p := i
				for {
					if in_path[p] != -1 {
						cycle := seq[in_path[p]:]
						for _, u := range cycle {
							sec := pref[u][next_val[u][head[u]]]
							nxt := pref[sec][tail[sec]]

							remove(nxt, rank[nxt][sec])
							if head[nxt] == -1 {
								fmt.Println("-1")
								return
							}

							remove(sec, rank[sec][nxt])
							if head[sec] == -1 {
								fmt.Println("-1")
								return
							}

							match[sec] = pref[sec][tail[sec]]
							Q = append(Q, nxt)
						}
						for _, node := range seq {
							in_path[node] = -1
						}
						break
					}
					in_path[p] = len(seq)
					seq = append(seq, p)

					sec := pref[p][next_val[p][head[p]]]
					nxt := pref[sec][tail[sec]]
					p = nxt
				}
				break
			}
		}

		if !found_cycle {
			break
		}

		if !processQ() {
			fmt.Println("-1")
			return
		}
	}

	for i := 1; i <= n; i++ {
		if head[i] == -1 || head[i] != tail[i] {
			fmt.Println("-1")
			return
		}
	}

	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = pref[i][head[i]]
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Print(" ")
		}
		fmt.Print(ans[i])
	}
	fmt.Println()
}
